// +build !solution

// Leave an empty line above this comment.
package chzap

import (
	"fmt"
	"strings"
	"time"
)

const timeFormat = "2006/01/02, 15:04:05"
const dateFormat = "2006/01/02"
const timeOnly = "15:04:05"
const timeLen = len(timeFormat)

type StatusChange struct {
	// Exported or Unexported?
	Time   time.Time
	IP     string
	Status string
}

// ChZap represent a channel change on a set-top box
type ChZap struct {
	// Exported or Unexported?
	Time     time.Time
	IP       string
	ToChan   string
	FromChan string
}

// NewSTBEvent creates a new set-top box(STB) event which can be either a ChZap or StatusChange
func NewSTBEvent(event string) (*ChZap, *StatusChange, error) {
	// ChZap format: {"2010/12/22, 20:22:32, 10.213.223.232, NRK2, NRK1", "20:22:32"}
	// StatusChange format: "2013/07/20, 21:57:42, 203.124.29.72, Volume: 50"}

	var eventSlice = strings.Split(event, ",")
	for i := 0; i < len(eventSlice); i++ {
		eventSlice[i] = strings.TrimSpace(eventSlice[i])
	}

	switch len(eventSlice) {
	case 5: // ChZap
		ip, toChan, fromChan := eventSlice[2], eventSlice[3], eventSlice[4]
		time, err := time.Parse(timeFormat, eventSlice[0]+", "+eventSlice[1])
		if err != nil {
			err = fmt.Errorf("NewSTBEvent: failed to parse timestamp")
			return nil, nil, err
		}
		chZap := ChZap{time, ip, toChan, fromChan}
		return &chZap, nil, nil
	case 4: // Statuschange
		ip, status := eventSlice[2], eventSlice[3]
		time, err := time.Parse(timeFormat, eventSlice[0]+", "+eventSlice[1])
		if err != nil {
			err = fmt.Errorf("NewSTBEvent: failed to parse timestamp")
			return nil, nil, err
		}
		staCha := StatusChange{time, ip, status}
		return nil, &staCha, nil
	case 3: // Error
		err := fmt.Errorf("NewSTBEvent: event with too few fields: %s, %s, %s ", eventSlice[0], eventSlice[1], eventSlice[2])
		return nil, nil, err
	case 2: // Error
		err := fmt.Errorf("NewSTBEvent: too short event string: %s, %s", eventSlice[0], eventSlice[1])
		return nil, nil, err
	}
	return nil, nil, nil
}

func (zap ChZap) String() string {
	return fmt.Sprintf("%s %s %s %s", zap.Time, zap.IP, zap.ToChan, zap.FromChan)
}

func (schg StatusChange) String() string {
	return fmt.Sprintf("%s %s", schg.Time, schg.Status)
}

// Duration returns between two zap events: The receiving (this) zap event and the provided event.
func (zap ChZap) Duration(provided time.Time) time.Duration {
	return zap.Time.Sub(provided)
}

// Date returns date of zap event in desired date format
func (zap ChZap) Date() string {
	return zap.Time.Format(dateFormat)
}
