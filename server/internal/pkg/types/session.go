package types

import "time"

type Session struct {
	ID        string    `json:"id"`
	UserID    uint32    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}
