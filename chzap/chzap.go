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

// StatusChange represents a change in status on a set-top box.
type StatusChange struct {
	// Exported or Unexported?
	Time       time.Time
	Volume     string
	MuteStatus string
	HDMIStatus string
}

// ChZap represent a channel change on a set-top box
type ChZap struct {
	// Exported or Unexported?
	Time     time.Time
	IP       string
	ToChan   string
	FromChan string
}

// Do we need to worry about wrong inputs? Wrong input should result in a error.

// NewSTBEvent creates a new set-top box(STB) event which can be either a ChZap or StatusChange
func NewSTBEvent(event string) (*ChZap, *StatusChange, error) {
	// ChZap format: {"2010/12/22, 20:22:32, 10.213.223.232, NRK2, NRK1", "20:22:32"}
	// StatusChange format: "2013/07/20, 21:57:42, 203.124.29.72, Volume: 50"}

	eventeSlice := strings.Split(event, ",")
	switch len(eventeSlice) {
	case 5: // ChZap
		eventTime, err := time.Parse(timeFormat, eventeSlice[0]+","+eventeSlice[1])
		if err != nil {
			err = fmt.Errorf("NewSTBEvent: failed to parse timestamp")
			return nil, nil, err
		}
		chZap := ChZap{Time: eventTime, IP: eventeSlice[2], ToChan: eventeSlice[3], FromChan: eventeSlice[4]}
		return &chZap, nil, err
	case 4: // Statuschange
		eventTime, err := time.Parse(timeFormat, eventeSlice[0]+","+eventeSlice[1])
		if err != nil {
			err = fmt.Errorf("NewSTBEvent: failed to parse timestamp")
			return nil, nil, err
		}
		staCha := StatusChange{Time: eventTime, Volume: eventeSlice[1], MuteStatus: eventeSlice[2], HDMIStatus: eventeSlice[3]}
		return nil, &staCha, nil
	case 3: // Error
		err := fmt.Errorf("NewSTBEvent: event with too few fields: %s,%s,%s", eventeSlice[0], eventeSlice[1], eventeSlice[2])
		return nil, nil, err
	case 2: // Error
		err := fmt.Errorf("NewSTBEvent: too short event string: %s,%s", eventeSlice[0], eventeSlice[1])
		return nil, nil, err
	default:
		// What is default case?
		// Maybe nil, nil, err: Unknown error ?
		// Or nil, nil, nil ?
		// Or no default case since nil, nil, nil is return if no case is triggered
	}
	return nil, nil, nil
}

func (zap ChZap) String() string {
	return fmt.Sprintf("%s%s%s%s", zap.Time, zap.IP, zap.ToChan, zap.FromChan)
}

func (schg StatusChange) String() string {
	return fmt.Sprintf("%s%s%s%s", schg.Time, schg.Volume, schg.MuteStatus, schg.HDMIStatus)
}

// Duration returns between two zap events: The receiving (this) zap event and the provided event.
func (zap ChZap) Duration(provided ChZap) time.Duration {
	return zap.Time.Sub(provided.Time)
}

// Date returns date of zap event in desired date format
func (zap ChZap) Date() string {
	return zap.Time.Format("2006/02/01")
}
