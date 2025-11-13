package models

import "time"

type LinkRequest struct {
	Links []string `json:"links"`
}

type LinkStatus struct {
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CheckTime time.Time `json:"check_time,omitempty"`
}

type Task struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	Status    string       `json:"status"`
	Results   []LinkStatus `json:"results"`
}
