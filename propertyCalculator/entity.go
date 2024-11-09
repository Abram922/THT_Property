package propertyCalculator

import (
	"time"
)

type Property struct {
	Timestamp  time.Time
	Type       string
	Area       float32
	BuildYear  int
	Location   string
	Corner     string
	Parking    int
	Facilities []string
}

var properties []Property
