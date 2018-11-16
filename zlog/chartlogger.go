package zlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

type Chartlogger struct {
	v    map[string][]*ViewTime // Key: Channelname, value: viewtime
	lock sync.Mutex
}
type ViewTime struct {
	Times time.Time
	Views float64
}

func NewChartLogger() *Chartlogger {
	cl := Chartlogger{v: make(map[string][]*ViewTime, 0)}
	return &cl
}

// LogZap updates count for the two channels in the zap
func (cl *Chartlogger) LogZap(z chzap.ChZap) {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	// Log views
	vtSlice, _ := cl.v[z.ToChan]
	if len(vtSlice) <= 0 {
		vtSlice = append(vtSlice, &ViewTime{Times: z.Time, Views: 1})
	} else {
		prevVt := vtSlice[len(vtSlice)-1]
		vtSlice = append(vtSlice, &ViewTime{Times: z.Time, Views: prevVt.Views + 1})
	}

	vtSlice = cl.v[z.FromChan]
	if len(vtSlice) <= 0 {
		vtSlice = append(vtSlice, &ViewTime{Times: z.Time, Views: -1})
	} else {
		prevVt := vtSlice[len(vtSlice)-1]
		vtSlice = append(vtSlice, &ViewTime{Times: z.Time, Views: prevVt.Views - 1})
	}
	for _, value := range vtSlice {
		fmt.Printf("Times: %v, Views: %v\n", value.Times, value.Views)
	}
	fmt.Println()
}

func (cl *Chartlogger) GetChartVal(channelName string) []*ViewTime {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.v[channelName]
}
