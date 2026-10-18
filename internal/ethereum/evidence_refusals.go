package ethereum

import (
	"github.com/WildanFrananda/kinetix-notary-service/internal/contractgen/evidence"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

var evidenceRefusals = []refusal{
	{evidence.EvidenceRegistryBatchAlreadyAnchoredErrorID(), "BatchAlreadyAnchored", notary.KindAlreadyAnchored},
	{evidence.EvidenceRegistryCloseInFutureErrorID(), "CloseInFuture", notary.KindAnchorDeferred},
	{
		evidence.EvidenceRegistryAccessControlUnauthorizedAccountErrorID(),
		"AccessControlUnauthorizedAccount",
		notary.KindNotAnchorer,
	},
	{evidence.EvidenceRegistryBatchConflictErrorID(), "BatchConflict", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryPreviousDigestMismatchErrorID(), "PreviousDigestMismatch", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryBatchOutOfSequenceErrorID(), "BatchOutOfSequence", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryRegistryIsSealedErrorID(), "RegistryIsSealed", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryBatchBeforeGenesisErrorID(), "BatchBeforeGenesis", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryUnknownStreamErrorID(), "UnknownStream", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryEmptyBatchErrorID(), "EmptyBatch", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryTooManyLeavesErrorID(), "TooManyLeaves", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryCloseNotAfterPreviousErrorID(), "CloseNotAfterPrevious", notary.KindAnchorRefused},
	{evidence.EvidenceRegistryPositionNotAfterPreviousErrorID(), "PositionNotAfterPrevious", notary.KindAnchorRefused},
	{
		evidence.EvidenceRegistrySafeCastOverflowedUintDowncastErrorID(),
		"SafeCastOverflowedUintDowncast",
		notary.KindAnchorRefused,
	},
}
