package entities

import "github.com/google/uuid"

type Topic struct {
	ID      uuid.UUID
	Name    string
	Pattern string
}

func NewTopic(id uuid.UUID, name string, pattern string) *Topic {
	return &Topic{
		ID:      id,
		Name:    name,
		Pattern: pattern,
	}
}
