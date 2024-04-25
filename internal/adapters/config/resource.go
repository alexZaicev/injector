package config

import "github.com/google/uuid"

type ResourceKind string

const (
	ResourceKindQueue ResourceKind = "Queue"
	ResourceKindTopic ResourceKind = "Topic"
)

type Resource struct {
	ApiVersion string         `yaml:"apiVersion,omitempty"`
	Kind       ResourceKind   `yaml:"kind,omitempty"`
	Metadata   map[string]any `yaml:"metadata,omitempty"`
	Spec       map[string]any `yaml:"spec,omitempty"`
}

type QueueSpec struct {
	ID   uuid.UUID `yaml:"-"`
	Name string    `yaml:"name,omitempty"`
}

type TopicSpec struct {
	ID      uuid.UUID `yaml:"-"`
	Name    string    `yaml:"name,omitempty"`
	Pattern string    `yaml:"pattern,omitempty"`
}
