// +build !solution

// Leave an empty line above this comment.
package zlog

import (
	"fmt"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

type Zaps []chzap.ChZap

func NewSimpleZapLogger() ZapLogger {
	zs := make(Zaps, 0)
	return &zs
}

func (zs *Zaps) LogZap(z chzap.ChZap) {
	*zs = append(*zs, z)
}

func (zs *Zaps) LogStatus(z chzap.StatusChange) {
}

func (zs *Zaps) Entries() int {
	return len(*zs)
}

func (zs *Zaps) String() string {
	return fmt.Sprintf("SS: %d", len(*zs))
}

// Viewers returns the current number of viewers for a channel.
func (zs *Zaps) Viewers(chName string) int {
	defer util.TimeElapsed(time.Now(), "simple.Viewers")
	viewers := 0

	for _, v := range *zs {
		if v.ToChan == chName {
			viewers++
		}
		if v.FromChan == chName {
			viewers--
		}
	}
	return viewers
}

// Channels creates a slice of the (unique) channels found in the zaps(both to and from).
func (zs *Zaps) Channels() []string {
	defer util.TimeElapsed(time.Now(), "simple.Channels")

	m := make(map[string]bool)    // Key: Channelname. To prevent duplicates.
	for _, channel := range *zs { // Create copy of zaps slice and range trough
		_, exists := m[channel.ToChan]
		if exists == false {
			m[channel.ToChan] = true
		}

		_, exists = m[channel.FromChan]
		if exists == false {
			m[channel.FromChan] = true
		}
	}

	// Create slice with unique channelnames
	channels := make([]string, 0)
	for channel := range m {
		channels = append(channels, channel)
	}

	return channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (zs *Zaps) ChannelsViewers() []*ChannelViewers {
	defer util.TimeElapsed(time.Now(), "simple.ChannelsViewers")

	m := make(map[string]int) // Key: Channelname. Value: Viewcount
	for _, channel := range *zs {
		count, exists := m[channel.ToChan]
		if exists {
			m[channel.ToChan] = count + 1
		} else {
			m[channel.ToChan] = 1
		}

		count, exists = m[channel.FromChan]
		if exists {
			m[channel.FromChan] = count - 1
		} else {
			m[channel.FromChan] = -1
		}
	}

	// Create a []*ChannelViewers slice from the map
	res := make([]*ChannelViewers, 0)
	for channel, viewers := range m {
		channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
		res = append(res, &channelViewer)
	}
	return res
}

func (zs *Zaps) StupidChart() ([]float64, []time.Time) {
	return nil, nil
}

func (zs *Zaps) ChartStats(views float64) {

}
