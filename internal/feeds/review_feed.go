package feeds

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/WildanFrananda/kinetix-notary-service/internal/notary"
	reviewv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/review/v1"
)

const reviewPeer = "review"

type ReviewFeed struct {
	conn    *grpc.ClientConn
	stub    reviewv1.ReviewEvidenceServiceClient
	timeout time.Duration
}

var _ notary.EvidenceFeed = (*ReviewFeed)(nil)

func DialReviews(settings Settings, options ...grpc.DialOption) (*ReviewFeed, error) {
	conn, err := dial("feeds.DialReviews", reviewPeer, settings, options)
	if err != nil {
		return nil, err
	}

	return &ReviewFeed{
		conn:    conn,
		stub:    reviewv1.NewReviewEvidenceServiceClient(conn),
		timeout: settings.Timeout,
	}, nil
}

func (f *ReviewFeed) Close() error {
	return closeConnection("feeds.ReviewFeed.Close", f.conn)
}

func (*ReviewFeed) Stream() notary.EvidenceStream {
	return notary.StreamReview
}

func (f *ReviewFeed) Since(ctx context.Context, after notary.LogPosition, limit int32) (notary.EvidencePage, error) {
	const op = "feeds.ReviewFeed.Since"

	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	response, err := f.stub.ReviewsSince(ctx, &reviewv1.ReviewsSinceRequest{After: cursorOf(after), Limit: limit})
	if err != nil {
		return notary.EvidencePage{}, classify(op, reviewPeer, err)
	}

	records := make([]notary.EvidenceRecord, 0, len(response.GetReviews()))
	for _, review := range response.GetReviews() {
		record, err := reviewRecordOf(op, review)
		if err != nil {
			return notary.EvidencePage{}, err
		}

		records = append(records, record)
	}

	return pageOf(op, reviewPeer, records, response.GetNext(), response.GetHasMore())
}

func reviewRecordOf(op string, review *reviewv1.ReviewEvidence) (notary.ReviewRecord, error) {
	at, err := positionOf(op, reviewPeer, "a review", review.GetPosition())
	if err != nil {
		return notary.ReviewRecord{}, err
	}

	if err := texts(
		op,
		reviewPeer,
		field{name: "a review id", value: review.GetReviewId()},
		field{name: "an order", value: review.GetOrderId()},
		field{name: "a subject", value: review.GetSubjectId()},
		field{name: "an author", value: review.GetAuthorPrincipalId()},
		field{name: "a merchant", value: review.GetMerchantPrincipalId(), optional: true},
	); err != nil {
		return notary.ReviewRecord{}, err
	}

	kind, err := enumOf(op, reviewPeer, "review kind", int32(review.GetKind()), true)
	if err != nil {
		return notary.ReviewRecord{}, err
	}

	rating, err := enumOf(op, reviewPeer, "rating", review.GetRating(), true)
	if err != nil {
		return notary.ReviewRecord{}, err
	}

	comment, err := hashOf(op, reviewPeer, "comment hash", review.GetCommentSha256(), false)
	if err != nil {
		return notary.ReviewRecord{}, err
	}

	createdAt, err := timeOf(op, reviewPeer, "a creation time", review.GetCreatedAt())
	if err != nil {
		return notary.ReviewRecord{}, err
	}

	return notary.ReviewRecord{
		At:                  at,
		ReviewID:            review.GetReviewId(),
		Kind:                kind,
		OrderID:             review.GetOrderId(),
		SubjectID:           review.GetSubjectId(),
		AuthorPrincipalID:   review.GetAuthorPrincipalId(),
		MerchantPrincipalID: review.GetMerchantPrincipalId(),
		Rating:              rating,
		CommentSHA256:       comment,
		CreatedAt:           createdAt,
	}, nil
}
