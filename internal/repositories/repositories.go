package repositories

import (
	"slices"
	"sync"
	"time"
	"todo-golang/internal/models"
)

const maxStoredEvents = 5000

var (
	eventsMu sync.RWMutex
	events   []models.Event
)

func StoreEvent(event models.Event) {
	eventsMu.Lock()
	defer eventsMu.Unlock()

	events = append(events, event)
	if len(events) > maxStoredEvents {
		events = slices.Clone(events[len(events)-maxStoredEvents:])
	}
}

func GetEventsSince(since time.Time) []models.Event {
	eventsMu.RLock()
	defer eventsMu.RUnlock()

	filtered := make([]models.Event, 0, len(events))
	for _, event := range events {
		if event.Timestamp.After(since) || event.Timestamp.Equal(since) {
			filtered = append(filtered, event)
		}
	}

	return filtered
}
