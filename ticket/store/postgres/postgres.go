package postgres

import (
	"context"
	"errors"
	"ticketAPI/commonerrors"
	"ticketAPI/ticket"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PsqlStore struct {
	conpool *pgxpool.Pool
}

func New(connStr string) (*PsqlStore, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &PsqlStore{
		conpool: pool,
	}, nil
}

func (p *PsqlStore) Close() {
	p.conpool.Close()
}

func (p *PsqlStore) GetTicket(ctx context.Context, id int) (ticket.Ticket, error) {
	sqlStatement := `SELECT * FROM tickets
	WHERE ticket_id = $1`

	t := ticket.Ticket{}
	row := p.conpool.QueryRow(ctx, sqlStatement, id)
	err := row.Scan(&t.ID, &t.Name, &t.Desc, &t.Allocation)
	if err != nil {
		return ticket.Ticket{}, errors.New(commonerrors.ErrNotFound)
	}
	return t, nil
}

func (p *PsqlStore) InsertTicketOption(ctx context.Context, t ticket.Ticket) (ticket.Ticket, error) {
	conn, err := p.conpool.Begin(ctx)
	if err != nil {
		return ticket.Ticket{}, err
	}

	sqlStatement := `INSERT INTO tickets (name, descr, allocation) 
	VALUES  ($1, $2, $3) 
	RETURNING ticket_id`

	row := conn.QueryRow(ctx, sqlStatement, t.Name, t.Desc, t.Allocation)
	if err != nil {
		return ticket.Ticket{}, err
	}

	err = row.Scan(&t.ID)
	if err != nil {
		conn.Rollback(ctx)
		return ticket.Ticket{}, err
	}

	err = conn.Commit(ctx)
	if err != nil {
		conn.Rollback(ctx)
		return ticket.Ticket{}, err
	}
	return t, nil
}

func (p *PsqlStore) PurchaseTicket(ctx context.Context, id int, alloc int) error {
	sqlStatement := `UPDATE tickets
	SET allocation = $1
	WHERE ticket_id=$2`

	_, err := p.conpool.Exec(ctx, sqlStatement, alloc, id)
	if err != nil {
		return err
	}
	return nil
}
