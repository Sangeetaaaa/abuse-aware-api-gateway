package queue

import "todo-golang/internal/models"

var EventChannel = make(chan models.Event, 1000)
