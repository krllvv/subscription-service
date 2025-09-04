package model

import (
	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	UserID    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   *string   `json:"end_date,omitempty"`
}

type SubRequest struct {
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	UserID    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   *string   `json:"end_date,omitempty"`
}
