package evidencing

import (
	"context"
	"log/slog"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
)

type Collector struct {
	scheme   notary.EvidenceScheme
	salts    notary.SaltSource
	feeds    []notary.EvidenceFeed
	store    notary.EvidenceStore
	observer Observer
	log      *slog.Logger
	policy   Policy
}

func NewCollector(
	scheme notary.EvidenceScheme,
	salts notary.SaltSource,
	feeds []notary.EvidenceFeed,
	store notary.EvidenceStore,
	observer Observer,
	log *slog.Logger,
	policy Policy,
) *Collector {
	return &Collector{
		scheme:   scheme,
		salts:    salts,
		feeds:    feeds,
		store:    store,
		observer: observer,
		log:      log,
		policy:   policy,
	}
}

func (c *Collector) Collect(ctx context.Context, lease notary.Lease) error {
	for _, feed := range c.feeds {
		err := c.drain(ctx, lease, feed)
		c.observer.EvidenceFeed(feed.Stream(), err == nil)

		if err == nil {
			continue
		}

		if fatal(err) {
			return err
		}

		logFailure(ctx, c.log, "evidence collection", feed.Stream(), err)
	}

	return nil
}

func (c *Collector) drain(ctx context.Context, lease notary.Lease, feed notary.EvidenceFeed) error {
	stream := feed.Stream()

	cursor, err := c.store.Cursor(ctx, stream)
	if err != nil {
		return err
	}

	for range c.policy.PagesPerTick {
		page, err := feed.Since(ctx, cursor, c.policy.PageSize)
		if err != nil {
			return err
		}

		if err := page.Follows(stream, cursor, c.policy.PageSize); err != nil {
			return err
		}

		entries := make([]notary.EvidenceEntry, len(page.Records))
		for i, record := range page.Records {
			entry, err := notary.EnterEvidence(c.scheme, c.salts, record)
			if err != nil {
				return err
			}

			entries[i] = entry
		}

		if err := c.store.Collect(ctx, lease, stream, cursor, entries, page.Next); err != nil {
			return err
		}

		if len(entries) > 0 {
			c.observer.EvidenceCollected(stream, len(entries))
			c.log.InfoContext(
				ctx,
				"evidence collected",
				"stream",
				stream.String(),
				"records",
				len(entries),
				"through_transaction",
				page.Next.Transaction,
				"through_sequence",
				page.Next.Sequence,
			)
		}

		cursor = page.Next
		if !page.HasMore {
			return nil
		}
	}

	return nil
}
