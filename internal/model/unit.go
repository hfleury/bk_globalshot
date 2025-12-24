package model

import "time"

type UnitType string

const (
	UnitTypeHouse UnitType = "HOUSE"
	UnitTypeFlat  UnitType = "FLAT"
)

type Unit struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      UnitType  `json:"type"`
	SiteID    string    `json:"site_id"`
	ClientID  *string   `json:"client_id,omitempty"` // Nullable
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
