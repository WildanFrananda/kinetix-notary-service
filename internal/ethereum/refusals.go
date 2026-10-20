package ethereum

import (
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/council"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledger"
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/ledgertwo"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var refusals = []refusal{
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
	{ledgertwo.LiabilityLedgerV2EpochNotProvenErrorID(), "EpochNotProven", notary.KindAnchorDeferred},
	{ledgertwo.LiabilityLedgerV2TotalNotHiddenErrorID(), "TotalNotHidden", notary.KindAnchorRefused},
}

var proofRefusals = []refusal{
	{ledgertwo.LiabilityLedgerV2StatementAlreadyProvenErrorID(), "StatementAlreadyProven", notary.KindAlreadyAnchored},
	{ledgertwo.LiabilityLedgerV2ProofRejectedErrorID(), "ProofRejected", notary.KindAnchorRefused},
	{ledgertwo.LiabilityLedgerV2EpochClosedErrorID(), "EpochClosed", notary.KindAnchorRefused},
	{ledgertwo.LiabilityLedgerV2LedgerIsSealedErrorID(), "LedgerIsSealed", notary.KindAnchorRefused},
	{ledgertwo.LiabilityLedgerV2TooManyLeavesErrorID(), "TooManyLeaves", notary.KindAnchorRefused},
	{ledgertwo.LiabilityLedgerV2NotAFieldElementErrorID(), "NotAFieldElement", notary.KindAnchorRefused},
}
