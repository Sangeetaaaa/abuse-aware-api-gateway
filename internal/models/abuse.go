package models

import "time"

type AbuseEvent struct {
	ID        string    `bson:"_id,omitempty"`
	Timestamp time.Time `bson:"timestamp"`
	IP        string    `bson:"ip"`
	Endpoint  string    `bson:"endpoint"`
	Action    string    `bson:"action"` // request, blocked, rate_limited
	Country   string    `bson:"country"`
	Score     int       `bson:"score"`
}
