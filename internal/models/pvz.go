package models

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID      uuid.UUID
	RegDate time.Time
	City    string
}
