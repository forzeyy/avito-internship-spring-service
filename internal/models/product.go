package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"date_time"`
	Type        string    `json:"type"`
	ReceptionID uuid.UUID `json:"reception_id"`
}
