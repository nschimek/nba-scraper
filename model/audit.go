package model

import "time"

type Audit struct {
	CreatedAt, UpdatedAt time.Time
}
