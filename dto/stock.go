package dto

import "time"

type Stock struct {
	ID        string     `json:"ID"`
	Name      string     `json:"Name"`
	CreatedAt *time.Time `json:"CreatedAt"`
	UpdatedAt *time.Time `json:"UpdatedAt"`
}
