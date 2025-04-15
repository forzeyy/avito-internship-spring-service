package models

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID      uuid.UUID `json:"id"`
	RegDate time.Time `json:"reg_date"`
	City    string    `json:"city"`
}
