package ethereum

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledgerinterface"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Ledger struct {
	client      *ethclient.Client
	binding     *ledger.LiabilityLedger
	contract    *bind.BoundContract
	expected    notary.LedgerIdentity
	interfaceID [4]byte
	scheme      uint16
	timeout     time.Duration
}

var _ notary.Ledger = (*Ledger)(nil)

func NewLedger(client *ethclient.Client, settings Settings, scheme notary.Scheme) *Ledger {
	binding := ledger.NewLiabilityLedger()
	return &Ledger{
		client:      client,
		binding:     binding,
		contract:    binding.Instance(client, common.Address(settings.Ledger.Address)),
		expected:    settings.Ledger,
		interfaceID: ledgerInterfaceID(),
		scheme:      scheme.Version(),
		timeout:     settings.Timeout,
	}
}

func (l *Ledger) Identity(ctx context.Context) (notary.LedgerIdentity, error) {
	const op = "ethereum.Ledger.Identity"
	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()

	chainID, err := l.client.ChainID(ctx)
	if err != nil {
		return notary.LedgerIdentity{}, unavailable(op, err)
	}
	if !chainID.IsUint64() || chainID.Uint64() != l.expected.ChainID {
		return notary.LedgerIdentity{}, notary.Errf(notary.KindWrongLedger, op, nil,
			"the RPC endpoint serves chain %s, the ledger is configured on chain %d", chainID, l.expected.ChainID)
	}
	code, err := l.client.CodeAt(ctx, common.Address(l.expected.Address), nil)
	if err != nil {
		return notary.LedgerIdentity{}, unavailable(op, err)
	}
	if len(code) == 0 {
		return notary.LedgerIdentity{}, notary.Errf(notary.KindWrongLedger, op, nil,
			"chain %d has no contract at %s", l.expected.ChainID, l.expected.Address.Hex())
	}
	supported, err := probe(l, ctx, l.binding.PackSupportsInterface(l.interfaceID), l.binding.UnpackSupportsInterface)
	if err != nil {
		return notary.LedgerIdentity{}, l.wrap(op, err)
	}
	if !supported {
		return notary.LedgerIdentity{}, notary.Errf(notary.KindWrongLedger, op, nil,
			"%s does not announce the liability ledger interface", l.expected.Address.Hex())
	}
	scheme, err := probe(l, ctx, l.binding.PackScheme(), l.binding.UnpackScheme)
	if err != nil {
		return notary.LedgerIdentity{}, l.wrap(op, err)
	}
	if scheme != l.scheme {
		return notary.LedgerIdentity{}, notary.Errf(notary.KindWrongLedger, op, nil,
			"%s commits with scheme %d, this verifier speaks scheme %d", l.expected.Address.Hex(), scheme, l.scheme)
	}
	return l.expected, nil
}

func (l *Ledger) Epoch(ctx context.Context, number uint64) (notary.AnchoredEpoch, error) {
	const op = "ethereum.Ledger.Epoch"
	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()

	final, err := l.finalizedBlock(ctx)
	if err != nil {
		return notary.AnchoredEpoch{}, unavailable(op, err)
	}
	read := func(block *big.Int) (ledger.Epoch, error) {
		return bind.Call(l.contract, &bind.CallOpts{Context: ctx, BlockNumber: block},
			l.binding.PackEpochAt(number), l.binding.UnpackEpochAt)
	}
	epoch, err := read(final)
	if err == nil {
		return l.anchoredEpoch(op, number, epoch)
	}
	if !l.isUnknownEpoch(err) && !errors.Is(err, bind.ErrNoCode) {
		return notary.AnchoredEpoch{}, l.classify(op, err)
	}
	_, err = read(nil)
	if err == nil {
		return notary.AnchoredEpoch{}, notary.Errf(notary.KindEpochNotFinal, op, nil,
			"epoch %d is anchored but not yet final at block %s", number, final)
	}
	return notary.AnchoredEpoch{}, l.classify(op, err)
}

func (l *Ledger) ConfirmsInclusion(ctx context.Context, claim notary.Claim) (bool, error) {
	const op = "ethereum.Ledger.ConfirmsInclusion"
	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()

	final, err := l.finalizedBlock(ctx)
	if err != nil {
		return false, unavailable(op, err)
	}
	siblings := make([]ledger.SumNode, len(claim.Proof.Siblings))
	for i, sibling := range claim.Proof.Siblings {
		siblings[i] = ledger.SumNode{Hash: sibling.Hash, Sum: sibling.Sum}
	}
	packed := l.binding.PackVerifyInclusion(claim.Epoch, claim.Entry.Commitment, claim.Entry.Amount,
		claim.Proof.Index, siblings)
	confirmed, err := bind.Call(l.contract, &bind.CallOpts{Context: ctx, BlockNumber: final}, packed,
		l.binding.UnpackVerifyInclusion)
	if err != nil {
		if l.isUnknownEpoch(err) || errors.Is(err, bind.ErrNoCode) {
			return false, notary.Errf(notary.KindEpochNotFinal, op, nil,
				"epoch %d is not final at block %s", claim.Epoch, final)
		}
		return false, l.classify(op, err)
	}
	return confirmed, nil
}

func (l *Ledger) finalizedBlock(ctx context.Context) (*big.Int, error) {
	header, err := l.client.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		return nil, fmt.Errorf("finalized block: %w", err)
	}
	return header.Number, nil
}

func (l *Ledger) isUnknownEpoch(err error) bool {
	raw, ok := revertData(err)
	if !ok {
		return false
	}
	return bytes.Equal(raw[:selectorLength], ledger.LiabilityLedgerUnknownEpochErrorID().Bytes()[:selectorLength])
}

func (l *Ledger) classify(op string, err error) error {
	if l.isUnknownEpoch(err) {
		raw, _ := revertData(err)
		unknown, decodeErr := l.binding.UnpackUnknownEpochError(raw[selectorLength:])
		if decodeErr == nil {
			return notary.Errf(notary.KindEpochNotAnchored, op, nil,
				"%s has no epoch %d", l.expected.Address.Hex(), unknown.Number)
		}
	}
	return unavailable(op, err)
}

func probe[T any](l *Ledger, ctx context.Context, calldata []byte, unpack func([]byte) (T, error)) (T, error) {
	var zero T
	address := common.Address(l.expected.Address)
	output, err := l.client.CallContract(ctx, goethereum.CallMsg{To: &address, Data: calldata}, nil)
	if err != nil {
		if _, reverted := revertData(err); reverted || isBareRevert(err) {
			return zero, errNotALedger
		}
		return zero, fmt.Errorf("call %s: %w", address.Hex(), err)
	}
	if len(output) != abiWordLength {
		return zero, errNotALedger
	}
	value, err := unpack(output)
	if err != nil {
		return zero, errNotALedger
	}
	return value, nil
}

func (l *Ledger) wrap(op string, err error) error {
	if errors.Is(err, errNotALedger) {
		return notary.Errf(notary.KindWrongLedger, op, nil,
			"%s does not answer as a liability ledger", l.expected.Address.Hex())
	}
	return unavailable(op, err)
}

func unavailable(op string, err error) error {
	return notary.Errf(notary.KindLedgerUnavailable, op, err, "")
}

func (l *Ledger) anchoredEpoch(op string, number uint64, epoch ledger.Epoch) (notary.AnchoredEpoch, error) {
	asOf, err := notary.UnixTimeFromSeconds(epoch.AsOf)
	if err != nil {
		return notary.AnchoredEpoch{}, notary.Errf(notary.KindWrongLedger, op, err, "epoch %d carries no valid snapshot time", number)
	}
	anchoredAt, err := notary.UnixTimeFromSeconds(epoch.AnchoredAt)
	if err != nil {
		return notary.AnchoredEpoch{}, notary.Errf(notary.KindWrongLedger, op, err, "epoch %d carries no valid anchor time", number)
	}
	return notary.AnchoredEpoch{
		Number:     number,
		Root:       epoch.Root,
		Digest:     epoch.Digest,
		Total:      epoch.Total,
		LeafCount:  epoch.LeafCount,
		AsOf:       asOf,
		AnchoredAt: anchoredAt,
	}, nil
}

func ledgerInterfaceID() [4]byte {
	var id [4]byte
	for _, method := range ledgerinterface.NewILiabilityLedger().GetABI().Methods {
		for i := range id {
			id[i] ^= method.ID[i]
		}
	}
	return id
}
