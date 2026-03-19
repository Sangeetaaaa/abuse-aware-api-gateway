package queue

import "todo-golang/internal/models"

var AbuseChannel = make(chan models.AbuseEvent, 1000)
