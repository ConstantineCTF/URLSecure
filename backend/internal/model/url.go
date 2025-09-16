package model

import "time"

// URL represents a shortened URL entry stored in the database.
type URL struct {
	ID        int       `db:"id"`                 // Primary key
	Code      string    `db:"code"`               // The unique short code
	Target    string    `db:"target"`             // The original (target) URL
	CreatedAt time.Time `db:"created_at"`         // Timestamp when shortened URL was created
	ExpiresAt time.Time `db:"expires_at,omitempty"` // Optional expiration time
	Clicks    int       `db:"clicks"`             // Number of times the link has been clicked
}
