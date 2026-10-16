package integration_test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/account7702"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/entrypoint"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/paymaster"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var canonicalEntryPoint = common.HexToAddress("0x4337084D9E255Ff0702461CF8895CE9E3b5Ff108")

type sponsorship struct {
	account   common.Address
	paymaster common.Address
}

func (c *chain) sponsorship(t *testing.T, courtAddress, withholdingAddress common.Address) sponsorship {
	t.Helper()

	ctx := context.Background()
	deployed, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(entrypoint.EntryPointMetaData.Bin),
		c.client,
		nil,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(ctx, c.client, tx.Hash())
	require.NoError(t, err)

	code, err := c.client.CodeAt(ctx, deployed, nil)
	require.NoError(t, err)
	require.NoError(t, c.rpc.CallContext(ctx, nil, "anvil_setCode", canonicalEntryPoint, hexBytes(code)))

	account, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(account7702.Simple7702AccountMetaData.Bin),
		c.client,
		nil,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(ctx, c.client, tx.Hash())
	require.NoError(t, err)

	input := paymaster.NewNotaryPaymaster().PackConstructor(
		c.from,
		big.NewInt(3600),
		canonicalEntryPoint,
		account,
		courtAddress,
		withholdingAddress,
		big.NewInt(10_000_000_000_000_000),
		86_400,
	)
	address, tx, err := bind.DeployContract(
		c.transactor(),
		common.FromHex(paymaster.NotaryPaymasterMetaData.Bin),
		c.client,
		input,
	)
	require.NoError(t, err)
	_, err = bind.WaitDeployed(ctx, c.client, tx.Hash())
	require.NoError(t, err)

	binding := entrypoint.NewEntryPoint()
	opts := c.transactor()
	opts.Value = big.NewInt(1_000_000_000_000_000_000)
	deposit, err := bind.Transact(binding.Instance(c.client, canonicalEntryPoint), opts, binding.PackDepositTo(address))
	require.NoError(t, err)
	_, err = bind.WaitMined(ctx, c.client, deposit.Hash())
	require.NoError(t, err)

	sponsor := paymaster.NewNotaryPaymaster()
	opts = c.transactor()
	opts.Value = big.NewInt(1_000_000_000_000_000_000)
	stake, err := bind.Transact(sponsor.Instance(c.client, address), opts, sponsor.PackAddStake())
	require.NoError(t, err)
	_, err = bind.WaitMined(ctx, c.client, stake.Hash())
	require.NoError(t, err)

	return sponsorship{account: account, paymaster: address}
}

func (c *chain) depositOf(t *testing.T, paymasterAddress common.Address) *big.Int {
	t.Helper()

	binding := entrypoint.NewEntryPoint()
	deposit, err := bind.Call(
		binding.Instance(c.client, canonicalEntryPoint),
		&bind.CallOpts{Context: context.Background()},
		binding.PackBalanceOf(paymasterAddress),
		binding.UnpackBalanceOf,
	)
	require.NoError(t, err)

	return deposit
}

func holderKeystore(t *testing.T) (*ecdsa.PrivateKey, string, string) {
	t.Helper()

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	dir := t.TempDir()
	encrypted, err := keystore.EncryptKey(
		&keystore.Key{Address: crypto.PubkeyToAddress(key.PublicKey), PrivateKey: key},
		"holder-password",
		keystore.LightScryptN,
		keystore.LightScryptP,
	)
	require.NoError(t, err)

	keystorePath := filepath.Join(dir, "holder.json")
	passwordPath := filepath.Join(dir, "holder.password")

	require.NoError(t, os.WriteFile(keystorePath, encrypted, 0o600))
	require.NoError(t, os.WriteFile(passwordPath, []byte("holder-password\n"), 0o600))

	return key, keystorePath, passwordPath
}

func registrationsDocument(
	t *testing.T,
	withholding notary.CourtIdentity,
	signer *ecdsa.PrivateKey,
	registration notary.Registration,
) string {
	t.Helper()

	digest := keccak.NewScheme().RegistrationDigest(withholding, registration)
	signature, err := crypto.Sign(digest[:], signer)
	require.NoError(t, err)

	signature[64] += 27

	document := map[string]any{
		"court": map[string]string{"chainId": "31337", "address": withholding.Address.Hex()},
		"epoch": "0",
		"registrations": []map[string]string{{
			"holder":    registration.Holder.Hex(),
			"account":   registration.Account.Hex(),
			"fromEpoch": "0",
			"kind":      notary.AccountCustomerBalance.String(),
			"digest":    digest.Hex(),
			"signature": notary.Signature(signature).Hex(),
		}},
	}
	encoded, err := json.Marshal(document)
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "registrations.json")
	require.NoError(t, os.WriteFile(path, encoded, 0o600))

	return path
}

func hexBytes(raw []byte) string {
	return "0x" + common.Bytes2Hex(raw)
}
