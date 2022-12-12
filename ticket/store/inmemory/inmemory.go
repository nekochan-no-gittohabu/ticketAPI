// This package is written for testing, it mocks postgresql database
package inmemory

import (
	"context"
	"errors"
	"sync"
	"ticketAPI/commonerrors"
	"ticketAPI/ticket"
)

type MockStore struct {
	m       map[int]ticket.Ticket
	counter int
	mu      sync.RWMutex
}

func New(data []ticket.Ticket) *MockStore {
	m := make(map[int]ticket.Ticket)
	for _, d := range data {
		m[d.ID] = d
	}

	return &MockStore{
		m:       m,
		counter: 1}
}

func (s *MockStore) GetTicket(ctx context.Context, id int) (ticket.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.m[id]
	if !ok {
		return ticket.Ticket{}, errors.New(commonerrors.ErrNotFound)
	}
	return t, nil
}
func (s *MockStore) InsertTicketOption(ctx context.Context, t ticket.Ticket) (ticket.Ticket, error) {
	s.mu.Lock()
	t.ID = s.counter
	s.m[s.counter] = t
	s.counter++
	s.mu.Unlock()
	return s.GetTicket(ctx, t.ID)
}

func (s *MockStore) PurchaseTicket(ctx context.Context, id int, alloc int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.m[id]
	t.Allocation = alloc
	s.m[id] = t

	return nil
}
