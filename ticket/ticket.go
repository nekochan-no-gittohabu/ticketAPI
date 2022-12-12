package ticket

import "github.com/google/uuid"

type Ticket struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Allocation int    `json:"allocation"`
}

type PurchaseRequest struct {
	Quantity int       `json:"quantity"`
	UserID   uuid.UUID `json:"user_id"`
}
