package ethereum

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/council"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var refusals = []struct {
	id   common.Hash
	name string
	kind notary.Kind
}{
	{ledger.LiabilityLedgerEpochAlreadyAnchoredErrorID(), "EpochAlreadyAnchored", notary.KindAlreadyAnchored},
	{ledger.LiabilityLedgerSnapshotInFutureErrorID(), "SnapshotInFuture", notary.KindAnchorDeferred},
	{
		ledger.LiabilityLedgerAccessControlUnauthorizedAccountErrorID(),
		"AccessControlUnauthorizedAccount",
		notary.KindNotAnchorer,
	},
	{ledger.LiabilityLedgerEpochConflictErrorID(), "EpochConflict", notary.KindAnchorRefused},
	{ledger.LiabilityLedgerPreviousDigestMismatchErrorID(), "PreviousDigestMismatch", notary.KindAnchorRefused},
	{ledger.LiabilityLedgerEpochOutOfSequenceErrorID(), "EpochOutOfSequence", notary.KindAnchorRefused},
	{ledger.LiabilityLedgerLedgerIsSealedErrorID(), "LedgerIsSealed", notary.KindAnchorRefused},
	{ledger.LiabilityLedgerEpochBeforeGenesisErrorID(), "EpochBeforeGenesis", notary.KindAnchorRefused},
	{ledger.LiabilityLedgerTooManyLeavesErrorID(), "TooManyLeaves", notary.KindAnchorRefused},
	{ledger.LiabilityLedgerEmptyEpochMismatchErrorID(), "EmptyEpochMismatch", notary.KindAnchorRefused},
	{ledger.LiabilityLedgerSnapshotNotAfterPreviousErrorID(), "SnapshotNotAfterPrevious", notary.KindAnchorRefused},
	{council.AnchorCouncilNotEnoughApprovalsErrorID(), "NotEnoughApprovals", notary.KindAnchorRefused},
	{council.AnchorCouncilNotAMemberErrorID(), "NotAMember", notary.KindAnchorRefused},
	{council.AnchorCouncilApprovalsOutOfOrderErrorID(), "ApprovalsOutOfOrder", notary.KindAnchorRefused},
}
