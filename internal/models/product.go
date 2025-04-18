package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	DateTime    time.Time
	Type        string
	ReceptionID uuid.UUID
}
