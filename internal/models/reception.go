package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID
	DateTime time.Time
	PVZID    uuid.UUID
	Status   string
	ClosedAt *time.Time
}
