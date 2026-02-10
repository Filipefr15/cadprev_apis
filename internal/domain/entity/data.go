package entity

import "time"

// Data represents the domain entity for external API data
// This is domain-independent and contains only business logic
type Data struct {
	ID          string
	Title       string
	Description string
	Value       float64
	Source      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewData creates a new Data entity with validation
func NewData(id, title, description, source string, value float64) *Data {
	now := time.Now()
	return &Data{
		ID:          id,
		Title:       title,
		Description: description,
		Value:       value,
		Source:      source,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the data entity
func (d *Data) Update(title, description string, value float64) {
	d.Title = title
	d.Description = description
	d.Value = value
	d.UpdatedAt = time.Now()
}
