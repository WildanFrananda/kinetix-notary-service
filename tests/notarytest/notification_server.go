package notarytest

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"google.golang.org/protobuf/proto"

	notificationv1 "github.com/WildanFrananda/kinetix-notary-service/internal/protogen/notification/v1"
)

type NotificationServer struct {
	notificationv1.UnimplementedNotificationServiceServer

	mu       sync.Mutex
	requests []*notificationv1.NotifyRequest
	callers  []string
	Answer   func(*notificationv1.NotifyRequest) (*notificationv1.NotifyResponse, error)
}

func (s *NotificationServer) Notify(
	ctx context.Context,
	request *notificationv1.NotifyRequest,
) (*notificationv1.NotifyResponse, error) {
	s.mu.Lock()
	s.requests = append(s.requests, proto.CloneOf(request))
	s.callers = append(s.callers, callerOf(ctx))
	count := len(s.requests)
	answer := s.Answer
	s.mu.Unlock()

	if answer != nil {
		return answer(request)
	}

	return &notificationv1.NotifyResponse{Accepted: true, NotificationId: fmt.Sprintf("NTF-%d", count)}, nil
}

func (s *NotificationServer) Requests() []*notificationv1.NotifyRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.requests)
}

func (s *NotificationServer) Callers() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.callers)
}
