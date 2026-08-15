package service

import (
	"errors"
	"fmt"

	"mailer/internal/model"
	"mailer/internal/store"
)

var ErrEmptyBody = errors.New("empty body")

type Service struct {
	store     *store.Store
	batchSize int
}

func New(s *store.Store, batchSize int) *Service {
	if batchSize <= 0 {
		batchSize = 1
	}
	return &Service{store: s, batchSize: batchSize}
}

func (svc *Service) Submit(id, to, body string, priority int) (*model.Message, error) {
	if body == "" {
		return nil, ErrEmptyBody
	}
	m := &model.Message{ID: id, To: to, Body: body, Priority: priority, Status: model.StatusQueued}
	if err := svc.store.Enqueue(m); err != nil {
		return nil, fmt.Errorf("submit %s: %w", id, err)
	}
	return m, nil
}

func (svc *Service) Get(id string) (*model.Message, error) {
	m, err := svc.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", id, err)
	}
	return m, nil
}

func (svc *Service) Pending() []*model.Message {
	return svc.store.Pending()
}

func (svc *Service) PendingIDs() []string {
	return svc.store.PendingIDs()
}

func (svc *Service) PrepareBatches() ([][]*model.Message, error) {
	msgs := svc.store.Pending()
	if len(msgs) == 0 {
		return nil, errors.New("no pending messages")
	}
	out := make([][]*model.Message, 0)
	for i := 0; i < len(msgs); i += svc.batchSize {
		end := i + svc.batchSize
		if end > len(msgs) {
			end = len(msgs)
		}
		out = append(out, msgs[i:end])
	}
	return out, nil
}

func (svc *Service) MarkSent(id string) error {
	if err := svc.store.MarkSent(id); err != nil {
		return fmt.Errorf("mark sent %s: %w", id, err)
	}
	return nil
}

func (svc *Service) MarkFailed(id string) error {
	if err := svc.store.MarkFailed(id); err != nil {
		return fmt.Errorf("mark failed %s: %w", id, err)
	}
	return nil
}
