package evidencing

import (
	"context"
	"log/slog"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

func logFailure(ctx context.Context, log *slog.Logger, step string, stream notary.EvidenceStream, err error) {
	kind := notary.KindOf(err)
	switch kind {
	case notary.KindAnchorDeferred:
		log.InfoContext(ctx, step+" deferred", "stream", stream.String(), "reason", err.Error())
	case notary.KindSourceUnavailable, notary.KindLedgerUnavailable, notary.KindStoreUnavailable, notary.KindLeaseLost:
		log.WarnContext(ctx, step+" will retry", "stream", stream.String(), "kind", kind.String(), "err", err.Error())
	default:
		log.ErrorContext(
			ctx,
			step+" stopped short",
			"stream",
			stream.String(),
			"kind",
			kind.String(),
			"err",
			err.Error(),
		)
	}
}

func fatal(err error) bool {
	switch notary.KindOf(err) {
	case notary.KindLeaseLost, notary.KindStoreUnavailable:
		return true
	default:
		return false
	}
}
