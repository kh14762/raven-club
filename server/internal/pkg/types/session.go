package types

import "time"

type Session struct {
	ID        string
	UserID    uint32
	ExpiresAt time.Time
}
