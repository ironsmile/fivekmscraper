package main

import "time"

type RunData struct {
	ID       int
	Age      uint8
	Name     string
	Place    string
	RunDate  time.Time
	Time     time.Duration
	Position uint32
	AvgSpeed float32 // kph
	Tempo    time.Duration
}
