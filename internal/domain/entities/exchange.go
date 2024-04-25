package entities

import "github.com/google/uuid"

type ExchangeKind int

const (
	ExchangeKindQueue ExchangeKind = iota + 1
	ExchangeKindTopic
)

type Exchange struct {
	ID      uuid.UUID
	Kind    ExchangeKind
	Name    string
	Pattern string
}

func NewQueueExchange(id uuid.UUID, name string) *Exchange {
	return &Exchange{
		ID:   id,
		Kind: ExchangeKindQueue,
		Name: name,
	}
}

func NewTopicExchange(id uuid.UUID, name string, pattern string) *Exchange {
	return &Exchange{
		ID:      id,
		Kind:    ExchangeKindTopic,
		Name:    name,
		Pattern: pattern,
	}
}
