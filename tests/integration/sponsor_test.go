package integration_test

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/config"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/withholding"
	"github.com/WildanFrananda/kinetix-notary-service/internal/ethereum"
	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsorcli"
	"github.com/WildanFrananda/kinetix-notary-service/internal/sponsoring"
	"github.com/WildanFrananda/kinetix-notary-service/internal/verifying"
	"github.com/WildanFrananda/kinetix-notary-service/internal/wiring"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

type sponsored struct {
	generations
	sponsorship
	bundler  *notarytest.LocalBundler
	keystore string
	password string
	holder   notary.Address
}

func sponsoredEstate(t *testing.T) sponsored {
	t.Helper()

	g := withDirectory(t, anchorEpochZero(t, 5, true))
	s := g.chain.sponsorship(t, g.court, g.withholding)
	bundlerKey, _ := g.chain.funded(t)
	key, keystorePath, passwordPath := holderKeystore(t)

	return sponsored{
		generations: g,
		sponsorship: s,
		bundler:     notarytest.StartLocalBundler(t, g.chain.client, bundlerKey, canonicalEntryPoint),
		keystore:    keystorePath,
		password:    passwordPath,
		holder:      notary.Address(crypto.PubkeyToAddress(key.PublicKey)),
	}
}

func (s sponsored) withholdingIdentity() notary.CourtIdentity {
	return notary.CourtIdentity{ChainID: s.identity.ChainID, Address: notary.Address(s.withholding)}
}

func (s sponsored) run(t *testing.T, args ...string) (int, []map[string]string) {
	t.Helper()

	env := map[string]string{
		config.RPCURLKey:           s.chain.url,
		config.ChainIDKey:          strconv.FormatUint(s.identity.ChainID, 10),
		config.DirectoryAddressKey: notary.Address(s.directory).Hex(),
		config.RPCTimeoutKey:       "10s",
		config.PaymasterAddressKey: notary.Address(s.paymaster).Hex(),
		config.BundlerURLKey:       s.bundler.URL,
		config.HolderKeystoreKey:   s.keystore,
		config.HolderUnlockFileKey: s.password,
		config.ReceiptWaitKey:      "20s",
	}
	container, err := wiring.Sponsor(func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	})
	require.NoError(t, err)

	var out bytes.Buffer

	code := -1
	err = container.Invoke(func(
		client *ethclient.Client,
		scheme notary.Scheme,
		statements notary.StatementScheme,
		directory notary.Directory,
		verifier *verifying.DirectoryVerifier,
		sponsor *sponsoring.Sponsor,
	) {
		defer client.Close()

		code = sponsorcli.Run(
			context.Background(),
			args,
			sponsorcli.Dependencies{
				Sponsor:    sponsor,
				Directory:  directory,
				Verifier:   verifier,
				Statements: statements,
				Scheme:     scheme.Version(),
			},
			&out,
		)
	})
	require.NoError(t, err)

	var reports []map[string]string

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		var report map[string]string
		require.NoError(t, json.Unmarshal(scanner.Bytes(), &report))

		reports = append(reports, report)
	}

	require.NotEmpty(t, reports)

	return code, reports
}

func (s sponsored) requestState(t *testing.T, registration notary.Registration, epoch uint64) notary.RequestState {
	t.Helper()

	digest := keccak.NewScheme().RegistrationDigest(s.withholdingIdentity(), registration)
	binding := withholding.NewWithholdingCourt()
	request, err := bind.Call(
		binding.Instance(s.chain.client, s.withholding),
		&bind.CallOpts{Context: context.Background()},
		binding.PackRequestOf(digest, epoch),
		binding.UnpackRequestOf,
	)
	require.NoError(t, err)

	return notary.RequestState(request.State)
}

func (s sponsored) balanceOfHolder(t *testing.T) *big.Int {
	t.Helper()

	balance, err := s.chain.client.BalanceAt(context.Background(), common.Address(s.holder), nil)
	require.NoError(t, err)

	return balance
}

func TestAHolderWithNoEtherRequestsAStatementThroughThePaymaster(t *testing.T) {
	s := sponsoredEstate(t)
	registration := notary.Registration{Holder: s.holder, Account: notary.Digest{0xac}, FromEpoch: 0}
	path := registrationsDocument(t, s.withholdingIdentity(), s.key, registration)
	before := s.chain.depositOf(t, s.paymaster)

	require.Zero(t, s.balanceOfHolder(t).Sign())

	code, reports := s.run(t, "request", path, "0")

	require.Equal(t, sponsorcli.ExitLanded, code, "%v", reports)
	require.Equal(t, "landed", reports[0]["verdict"])
	require.Equal(t, notary.RequestOpen, s.requestState(t, registration, 0))
	require.Zero(t, s.balanceOfHolder(t).Sign())

	delegated, err := s.chain.client.CodeAt(context.Background(), common.Address(s.holder), nil)
	require.NoError(t, err)
	require.Equal(t, append([]byte{0xef, 0x01, 0x00}, s.account.Bytes()...), delegated)

	paid := s.chain.depositOf(t, s.paymaster)
	require.Equal(t, -1, paid.Cmp(before))

	early, reports := s.run(t, "unanswered", path, "0")

	require.Equal(t, sponsorcli.ExitNothingToDo, early, "%v", reports)
	require.Equal(t, notary.KindNotDue.String(), reports[0]["verdict"])
	require.Equal(t, notary.RequestOpen, s.requestState(t, registration, 0))
	require.Zero(t, s.chain.depositOf(t, s.paymaster).Cmp(paid))

	ctx := context.Background()
	require.NoError(t, s.chain.rpc.CallContext(ctx, nil, "evm_increaseTime", 24*3600+60))
	require.NoError(t, s.chain.rpc.CallContext(ctx, nil, "evm_mine"))

	overdue, reports := s.run(t, "unanswered", path, "0")

	require.Equal(t, sponsorcli.ExitLanded, overdue, "%v", reports)
	require.Equal(t, "landed", reports[0]["verdict"])
	require.Equal(t, notary.RequestUnanswered, s.requestState(t, registration, 0))
	require.Zero(t, s.balanceOfHolder(t).Sign())
	require.Equal(t, -1, s.chain.depositOf(t, s.paymaster).Cmp(paid))
}

func TestACallTheCourtWouldRefuseIsNeverSentAndCostsTheDepositNothing(t *testing.T) {
	s := sponsoredEstate(t)
	registration := notary.Registration{Holder: s.holder, Account: notary.Digest{0xac}, FromEpoch: 0}
	path := registrationsDocument(t, s.withholdingIdentity(), s.key, registration)

	never, reports := s.run(t, "unanswered", path, "0")

	require.Equal(t, sponsorcli.ExitNothingToDo, never, "%v", reports)
	require.Equal(t, notary.KindCallRefused.String(), reports[0]["verdict"])
	require.Contains(t, reports[0]["reason"], "no open request for epoch 0")

	first, reports := s.run(t, "request", path, "0")
	require.Equal(t, sponsorcli.ExitLanded, first, "%v", reports)

	before := s.chain.depositOf(t, s.paymaster)
	again, reports := s.run(t, "request", path, "0")

	require.Equal(t, sponsorcli.ExitNothingToDo, again, "%v", reports)
	require.Equal(t, notary.KindCallRefused.String(), reports[0]["verdict"])
	require.Contains(t, reports[0]["reason"], "still open")
	require.Zero(t, s.chain.depositOf(t, s.paymaster).Cmp(before))
}

func TestARegistrationTheNotaryNeverSignedIsNeverSponsored(t *testing.T) {
	s := sponsoredEstate(t)
	registration := notary.Registration{Holder: s.holder, Account: notary.Digest{0xac}, FromEpoch: 0}
	forger, _ := notarytest.StatementKey(t, "integration-forger")
	path := registrationsDocument(t, s.withholdingIdentity(), forger, registration)
	before := s.chain.depositOf(t, s.paymaster)

	code, reports := s.run(t, "request", path, "0")

	require.Equal(t, sponsorcli.ExitNotSponsored, code, "%v", reports)
	require.Equal(t, notary.KindNotSponsored.String(), reports[0]["verdict"])
	require.Equal(t, notary.RequestNone, s.requestState(t, registration, 0))
	require.Zero(t, s.chain.depositOf(t, s.paymaster).Cmp(before))
}

func TestThePaymasterAdapterReadsTheDeployedPaymaster(t *testing.T) {
	g := withDirectory(t, anchorEpochZero(t, 4, true))
	s := g.chain.sponsorship(t, g.court, g.withholding)
	ctx := context.Background()

	adapter := ethereum.NewPaymaster(g.chain.client, g.identity.ChainID, notary.Address(s.paymaster), 10*time.Second)

	identity, err := adapter.Identity(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		notary.PaymasterIdentity{
			ChainID:     g.identity.ChainID,
			Address:     notary.Address(s.paymaster),
			EntryPoint:  notary.Address(canonicalEntryPoint),
			Account:     notary.Address(s.account),
			Court:       notary.Address(g.court),
			Withholding: notary.Address(g.withholding),
			Signer:      g.signer,
			Ceiling:     big.NewInt(10_000_000_000_000_000),
		},
		identity,
	)

	funds, err := adapter.Funds(ctx)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1_000_000_000_000_000_000), funds.Deposit)
	require.Equal(t, big.NewInt(1_000_000_000_000_000_000), funds.LockedStake())
	require.True(t, funds.Staked)

	elsewhere := ethereum.NewSponsoredCalls().Execute(notary.SponsoredCall{Target: notary.Address(g.address)})
	_, err = adapter.Sponsorship(ctx, notary.Address{0x01}, elsewhere, nil)
	require.Equal(t, notary.KindNotSponsored, notary.KindOf(err), "%v", err)
	require.ErrorContains(t, err, "sponsors no call to "+g.address.Hex())

	for name, address := range map[string]common.Address{
		"the ledger":           g.address,
		"an account with code": s.account,
		"an empty address":     {0x0e},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := ethereum.NewPaymaster(
				g.chain.client,
				g.identity.ChainID,
				notary.Address(address),
				10*time.Second,
			).Identity(ctx)

			require.Equal(t, notary.KindNotSponsored, notary.KindOf(err), "%v", err)
		})
	}

	_, err = ethereum.NewPaymaster(g.chain.client, 1, notary.Address(s.paymaster), 10*time.Second).Identity(ctx)
	require.Equal(t, notary.KindNotSponsored, notary.KindOf(err), "%v", err)
}

func (s sponsored) issued(
	t *testing.T,
	registration notary.Registration,
	signer *ecdsa.PrivateKey,
) notary.IssuedRegistration {
	t.Helper()

	digest := keccak.NewScheme().RegistrationDigest(s.withholdingIdentity(), registration)
	signature, err := crypto.Sign(digest[:], signer)
	require.NoError(t, err)

	signature[64] += 27

	return notary.IssuedRegistration{
		Registration: registration,
		Kind:         notary.AccountCustomerBalance,
		Digest:       digest,
		Signature:    notary.Signature(signature),
	}
}

func TestThePaymasterNamesWhyItWouldNotSponsorACall(t *testing.T) {
	s := sponsoredEstate(t)
	ctx := context.Background()
	adapter := ethereum.NewPaymaster(s.chain.client, s.identity.ChainID, notary.Address(s.paymaster), 10*time.Second)
	calls := ethereum.NewSponsoredCalls()
	withholdingAt := notary.Address(s.withholding)
	forger, _ := notarytest.StatementKey(t, "integration-forger")
	mine := notary.Registration{Holder: s.holder, Account: notary.Digest{0xac}, FromEpoch: 0}
	theirs := notary.Registration{Holder: notary.Address{0x0b}, Account: notary.Digest{0xac}, FromEpoch: 0}
	sound := s.sign(t, s.claim(t, 1), s.key, s.signer)

	for name, c := range map[string]struct {
		call notary.SponsoredCall
		kind notary.Kind
		says string
	}{
		"an epoch not anchored": {
			calls.RequestStatement(withholdingAt, s.issued(t, mine, s.key), 9),
			notary.KindEpochNotAnchored,
			"epoch 9 is not anchored",
		},
		"an epoch before the registration": {
			calls.RequestStatement(
				withholdingAt,
				s.issued(t, notary.Registration{Holder: s.holder, Account: notary.Digest{0xac}, FromEpoch: 3}, s.key),
				0,
			),
			notary.KindCallRefused,
			"comes before the registration",
		},
		"no open request": {
			calls.RecordUnanswered(withholdingAt, s.issued(t, mine, s.key), 0),
			notary.KindCallRefused,
			"no open request for epoch 0",
		},
		"a sound proof": {
			calls.ProveBadProof(
				notary.Address(s.court),
				notary.StatementOf(keccak.NewScheme(), sound),
				sound.Statement.Signature,
				sound.Proof,
			),
			notary.KindCallRefused,
			"holds against the anchored root",
		},
		"a registration the notary never signed": {
			calls.RequestStatement(withholdingAt, s.issued(t, mine, forger), 0),
			notary.KindNotSponsored,
			"the signature recovers to",
		},
		"another holder's registration": {
			calls.RequestStatement(withholdingAt, s.issued(t, theirs, s.key), 0),
			notary.KindNotSponsored,
			"the registration is " + common.Address(theirs.Holder).Hex() + "'s",
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := adapter.Sponsorship(ctx, s.holder, calls.Execute(c.call), c.call.PaymasterData)

			require.Equal(t, c.kind, notary.KindOf(err), "%v", err)
			require.ErrorContains(t, err, c.says)
		})
	}

	request := calls.RequestStatement(withholdingAt, s.issued(t, mine, s.key), 0)
	window, err := adapter.Sponsorship(ctx, s.holder, calls.Execute(request), nil)

	require.NoError(t, err)
	require.True(t, window.After.IsZero())
	require.False(t, window.Until.IsZero())
}
