package models

import "time"

type User struct {
	ID   uint64
	Name string
}

type Event struct {
	ID        uint64
	Name      string
	Timestamp time.Time
	UserIDs   []uint64
}
