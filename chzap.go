// +build !solution

// Leave an empty line above this comment.
package lab6

import (
	"time"
)

const timeFormat = "2006/01/02, 15:04:05"
const dateFormat = "2006/01/02"
const timeOnly = "15:04:05"
const timeLen = len(timeFormat)

type StatusChange struct {
	// Exported or Unexported?
	Volume     int
	MuteStatus int
	HDMIStatus int
	//TODO finish this struct (1p)
}

type ChZap struct {
	Date         time.Time // trenger vi denne?
	Time         time.Time
	IP           string
	ToChan       string
	FromChan     string
	StatusChange StatusChange // skal det være status change object
	//TODO finish this struct (1p)
}

func NewSTBEvent(event string) (*ChZap, *StatusChange, error) {
	//TODO write this method (5p)
	return nil, nil, nil
}

func (zap ChZap) String() string {
	//TODO write this method (2p)
	return ""
}

func (schg StatusChange) String() string {
	//TODO write this method (1p)
	return ""
}

// The duration between receiving (this) zap event and the provided event
func (zap ChZap) Duration(provided ChZap) time.Duration {
	//TODO write this method (1p)
	return time.Duration(0)
}

func (zap ChZap) Date() string {
	//TODO write this method (1p)
	return ""
}
