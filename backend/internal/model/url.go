package model

import "time"

type URL struct {
	ID        int       `db:"id"`
	Code      string    `db:"code"`
	Target    string    `db:"target"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at,omitempty"`
	Clicks    int       `db:"clicks"`
}
