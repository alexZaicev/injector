package broker

import (
	"sync"

	"github.com/google/uuid"
)

// ChannelMap provide concurrency safe read-writes to map containing consumer channels
type ChannelMap struct {
	channels map[uuid.UUID]*Channel
	mu       sync.RWMutex
}

func newChannelMap() *ChannelMap {
	return &ChannelMap{
		channels: make(map[uuid.UUID]*Channel),
	}
}

func (cm *ChannelMap) Add(channel *Channel) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.channels[channel.queue.ID] = channel
}

func (cm *ChannelMap) Remove(queueID uuid.UUID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.channels, queueID)
}

func (cm *ChannelMap) Get(queueID uuid.UUID) (*Channel, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	channel, ok := cm.channels[queueID]
	return channel, ok
}

func (cm *ChannelMap) Find(queueID uuid.UUID) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	_, ok := cm.channels[queueID]
	return ok
}

func (cm *ChannelMap) Keys() []uuid.UUID {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	keys := make([]uuid.UUID, 0, len(cm.channels))

	for queueID := range cm.channels {
		keys = append(keys, queueID)
	}

	return keys
}
