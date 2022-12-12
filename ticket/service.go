package ticket

import (
	"context"
	"errors"
	"ticketAPI/commonerrors"
)

type Store interface {
	GetTicket(ctx context.Context, id int) (Ticket, error)
	InsertTicketOption(ctx context.Context, t Ticket) (Ticket, error)
	PurchaseTicket(ctx context.Context, id int, alloc int) error
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetTicket(ctx context.Context, id int) (Ticket, error) {
	return s.store.GetTicket(ctx, id)
}

func (s *Service) CreateTicketOption(ctx context.Context, t Ticket) (Ticket, error) {
	if t.Allocation <= 0 {
		return Ticket{}, errors.New(commonerrors.ErrParameter)
	}
	return s.store.InsertTicketOption(ctx, t)
}

func (s *Service) PurchaseTicket(ctx context.Context, id int, r PurchaseRequest) error {
	t, err := s.store.GetTicket(ctx, id)
	if err != nil {
		return err
	}

	alloc := 0
	if r.Quantity <= 0 {
		return errors.New(commonerrors.ErrParameter)
	}
	if r.Quantity > t.Allocation {
		return errors.New(commonerrors.ErrTooManyPurchases)
	}
	if t.Allocation-r.Quantity > 0 {
		alloc = t.Allocation - r.Quantity
	}

	return s.store.PurchaseTicket(ctx, id, alloc)
}
