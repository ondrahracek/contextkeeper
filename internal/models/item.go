package models

// Item represents a context item
type Item struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Created string `json:"created"`
}
