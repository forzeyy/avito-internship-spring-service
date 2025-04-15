package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID  `json:"id"`
	DateTime time.Time  `json:"date_time"`
	PVZID    uuid.UUID  `json:"pvz_id"`
	Status   string     `json:"status"`
	ClosedAt *time.Time `json:"closed_at,omitempty"`
}
