package model

import "time"

type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UnitID    string    `json:"unit_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
