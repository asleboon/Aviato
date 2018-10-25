// +build !solution

// Leave an empty line above this comment.
package lab6

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
	Time       time.Time
	Volume     string
	MuteStatus string
	HDMIStatus string
	//TODO finish this struct (1p)
}

type ChZap struct {
	Time     time.Time
	IP       string
	ToChan   string
	FromChan string
	// Statuschange object?
	//TODO finish this struct (1p)
}

func NewSTBEvent(event string) (*ChZap, *StatusChange, error) {
	// TODO write this method (5p)
	// Input string format ChZap:
	//{"2010/12/22, 20:22:32, 10.213.223.232, NRK2, NRK1", "20:22:32"} len = 5

	// Input string format StatusChange:
	// "2013/07/20, 21:57:42, 203.124.29.72, Volume: 50"} len = 4

	// Split it into substrings
	eventStr := strings.Split(event, ",")
	// Teste om det er en statusChange eller en ChZap
	// Do we need to worry about wrong inputs?
	switch len(eventStr) {
	case 5:
		// it is a ChZap
		// Do we need error handeling?
		// format string to Datetime object
		eventTime, err := time.Parse(timeFormat, eventStr[0]+","+eventStr[1])
		// fmt.Printf("\nEventDate: %q\n", eventDate)
		if err != nil {
			err = fmt.Errorf("NewSTBEvent: failed to parse timestamp")
			return nil, nil, err
		}
		chZap := ChZap{Time: eventTime, IP: eventStr[2], ToChan: eventStr[3], FromChan: eventStr[4]}
		return &chZap, nil, err
	case 4:
		// it is a StatusChange
		eventTime, err := time.Parse(timeFormat, eventStr[0]+","+eventStr[1])
		fmt.Printf("\neventTime %q\n", eventTime)
		if err != nil {
			err = fmt.Errorf("NewSTBEvent: failed to parse timestamp")
			return nil, nil, err
		}
		staCha := StatusChange{Time: eventTime, Volume: eventStr[1], MuteStatus: eventStr[2], HDMIStatus: eventStr[3]}
		return nil, &staCha, nil
	case 3:
		// wrong input
		err := fmt.Errorf("NewSTBEvent: event with too few fields: %s,%s,%s", eventStr[0], eventStr[1], eventStr[2])
		return nil, nil, err
	case 2:
		// wrong input
		err := fmt.Errorf("NewSTBEvent: too short event string: %s,%s", eventStr[0], eventStr[1])
		return nil, nil, err
	default:

	}

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
