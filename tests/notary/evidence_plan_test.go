package notary_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/WildanFrananda/kinetix-notary-service/internal/keccak"
	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	"github.com/WildanFrananda/kinetix-notary-service/tests/notarytest"
)

func page(records []notary.EvidenceRecord, next notary.LogPosition, more bool) notary.EvidencePage {
	return notary.EvidencePage{Records: records, Next: next, HasMore: more}
}

func TestAPageFollowsOnlyWhenItContinuesTheLog(t *testing.T) {
	records := notarytest.Records(notarytest.Deliveries(t, 3))
	end := records[2].Position()

	require.NoError(t, page(records, end, false).Follows(notary.StreamDelivery, notary.LogPosition{}, 3))
	require.NoError(t, page(nil, end, false).Follows(notary.StreamDelivery, end, 3))

	for name, c := range map[string]struct {
		page  notary.EvidencePage
		after notary.LogPosition
		limit int32
	}{
		"more than asked":       {page(records, end, false), notary.LogPosition{}, 2},
		"no limit":              {page(records, end, false), notary.LogPosition{}, 0},
		"another stream":        {page(notarytest.Records(notarytest.Reviews(t, 1)), notarytest.Reviews(t, 1)[0].At, false), notary.LogPosition{}, 3},
		"not after the cursor":  {page(records, end, false), records[0].Position(), 3},
		"out of order":          {page([]notary.EvidenceRecord{records[1], records[0]}, records[0].Position(), false), notary.LogPosition{}, 3},
		"next is not the end":   {page(records, records[1].Position(), false), notary.LogPosition{}, 3},
		"empty yet more":        {page(nil, notary.LogPosition{}, true), notary.LogPosition{}, 3},
		"empty and moved ahead": {page(nil, end, false), notary.LogPosition{}, 3},
	} {
		err := c.page.Follows(notary.StreamDelivery, c.after, c.limit)
		require.Equal(t, notary.KindSourceRefused, notary.KindOf(err), name)
	}
}

func TestAnEntryCommitsToItsRecordWithAFreshSalt(t *testing.T) {
	scheme := keccak.NewEvidence()
	record := notarytest.Reviews(t, 1)[0]

	entry, err := notary.EnterEvidence(scheme, &notarytest.Salts{Seed: 1}, record)
	require.NoError(t, err)
	require.NoError(t, entry.Holds(scheme))

	digest, err := scheme.RecordDigest(record)
	require.NoError(t, err)
	require.Equal(t, digest, entry.RecordDigest)
	require.Equal(t, scheme.EvidenceLeaf(notary.StreamReview, entry.Salt, digest), entry.Leaf)

	_, err = notary.EnterEvidence(scheme, &notarytest.Salts{Zero: true}, record)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	changed := entry
	reviewed := record
	reviewed.Rating = 5
	changed.Record = reviewed
	require.Equal(t, notary.KindDisagreement, notary.KindOf(changed.Holds(scheme)))

	resalted := entry
	resalted.Salt = notary.Salt{0x01}
	require.Equal(t, notary.KindDisagreement, notary.KindOf(resalted.Holds(scheme)))
}

func TestABatchPlanBindsItsRecordsInFeedOrder(t *testing.T) {
	scheme := keccak.NewEvidence()
	salts := &notarytest.Salts{Seed: 3}
	records := notarytest.Records(notarytest.Deliveries(t, 5))
	entries := make([]notary.EvidenceEntry, len(records))

	for i, record := range records {
		entry, err := notary.EnterEvidence(scheme, salts, record)
		require.NoError(t, err)

		entries[i] = entry
	}

	closedAt, err := notary.UnixTimeFromSeconds(1_760_000_000)
	require.NoError(t, err)

	plan, err := notary.PlanEvidenceBatch(scheme, notary.StreamDelivery, 4, notary.Digest{0x04}, closedAt, entries)
	require.NoError(t, err)
	require.Equal(t, uint64(5), plan.Header.LeafCount)
	require.Equal(t, records[4].Position(), plan.Header.Through)
	require.Equal(t, scheme.BatchDigest(plan.Header), plan.Digest)
	require.Equal(t, uint16(1), plan.Scheme)

	swapped := append([]notary.EvidenceEntry{}, entries...)
	swapped[1], swapped[2] = swapped[2], swapped[1]
	_, err = notary.PlanEvidenceBatch(scheme, notary.StreamDelivery, 4, notary.Digest{}, closedAt, swapped)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	_, err = notary.PlanEvidenceBatch(scheme, notary.StreamReview, 4, notary.Digest{}, closedAt, entries)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	_, err = notary.PlanEvidenceBatch(scheme, notary.StreamDelivery, 4, notary.Digest{}, notary.UnixTime{}, entries)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	_, err = notary.PlanEvidenceBatch(scheme, notary.StreamDelivery, 4, notary.Digest{}, closedAt, nil)
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))

	tampered := append([]notary.EvidenceEntry{}, entries...)
	tampered[3].Leaf = notary.Digest{0xff}
	_, err = notary.PlanEvidenceBatch(scheme, notary.StreamDelivery, 4, notary.Digest{}, closedAt, tampered)
	require.Equal(t, notary.KindDisagreement, notary.KindOf(err))
}

func TestAStreamIsNamedTheWayADocumentNamesIt(t *testing.T) {
	for _, stream := range notary.EvidenceStreams() {
		named, err := notary.EvidenceStreamNamed(stream.String())
		require.NoError(t, err)
		require.Equal(t, stream, named)
	}

	_, err := notary.EvidenceStreamNamed("unknown")
	require.Equal(t, notary.KindInvalidArgument, notary.KindOf(err))
}
