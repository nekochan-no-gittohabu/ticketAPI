package ticket_test

import (
	"context"
	"errors"
	"testing"
	"ticketAPI/commonerrors"
	"ticketAPI/ticket"
	"ticketAPI/ticket/store/inmemory"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestGetTicket(t *testing.T) {
	store := inmemory.New(getDataSet())
	service := ticket.New(store)
	ctx := context.Background()

	tests := []struct {
		name string
		id   int
		want ticket.Ticket
		err  error
	}{
		{
			name: "ticket1",
			id:   getDataSet()[0].ID,
			want: getDataSet()[0],
			err:  nil,
		},
		{
			name: "ticket2",
			id:   getDataSet()[1].ID,
			want: getDataSet()[1],
			err:  nil,
		},
		{
			name: "ticket3",
			id:   getDataSet()[2].ID,
			want: getDataSet()[2],
			err:  nil,
		},
		{
			name: "not found",
			id:   45,
			want: ticket.Ticket{},
			err:  errors.New(commonerrors.ErrNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetTicket(ctx, tt.id)
			if !errors.Is(err, tt.err) {
				t.Fatal("failed to get ticket", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestCreateTicketOption(t *testing.T) {
	store := inmemory.New(nil)
	service := ticket.New(store)
	ctx := context.Background()

	tests := []struct {
		name string
		req  ticket.Ticket
		want ticket.Ticket
		err  bool
	}{
		{
			name: "ok",
			req: ticket.Ticket{
				Name:       "ticket1",
				Desc:       "desc1",
				Allocation: 3,
			},
			want: ticket.Ticket{
				ID:         1,
				Name:       "ticket1",
				Desc:       "desc1",
				Allocation: 3,
			},
			err: false,
		},
		{
			name: "ticket with ID",
			req: ticket.Ticket{
				ID:         75,
				Name:       "ticket1",
				Desc:       "desc1",
				Allocation: 30,
			},
			want: ticket.Ticket{
				ID:         2,
				Name:       "ticket1",
				Desc:       "desc1",
				Allocation: 30,
			},
			err: false,
		},
		{
			name: "invalid allocation",
			req: ticket.Ticket{
				ID:         75,
				Name:       "ticket1",
				Desc:       "desc1",
				Allocation: -30,
			},
			want: ticket.Ticket{},
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CreateTicketOption(ctx, tt.req)
			if (err != nil && !tt.err) || (err == nil && tt.err) {
				t.Fatal("failed to get ticket")
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestPurchaseTicket(t *testing.T) {
	store := inmemory.New(getDataSet())
	service := ticket.New(store)
	ctx := context.Background()

	tests := []struct {
		name string
		id   int
		req  ticket.PurchaseRequest
		err  bool
	}{
		{
			name: "ok",
			id:   getDataSet()[0].ID,
			req: ticket.PurchaseRequest{
				Quantity: 1,
				UserID:   uuid.New(),
			},
			err: false,
		},
		{
			name: "too many tickets",
			id:   getDataSet()[1].ID,
			req: ticket.PurchaseRequest{
				Quantity: 50,
				UserID:   uuid.New(),
			},
			err: true,
		},
		{
			name: "negative purchase",
			id:   getDataSet()[2].ID,
			req: ticket.PurchaseRequest{
				Quantity: -100,
				UserID:   uuid.New(),
			},
			err: true,
		},
		{
			name: "not found",
			id:   78,
			req: ticket.PurchaseRequest{
				Quantity: 100,
				UserID:   uuid.New(),
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.PurchaseTicket(ctx, tt.id, tt.req)
			if (err != nil && !tt.err) || (err == nil && tt.err) {
				t.Fatal("failed to get ticket")
			}
		})
	}
}

func getDataSet() []ticket.Ticket {
	return []ticket.Ticket{
		{
			ID:         1,
			Name:       "ticket1",
			Desc:       "desc1",
			Allocation: 19,
		},
		{
			ID:         2,
			Name:       "ticket2",
			Desc:       "desc2",
			Allocation: 32,
		},
		{
			ID:         3,
			Name:       "ticket3",
			Desc:       "desc3",
			Allocation: 0,
		},
	}

}
