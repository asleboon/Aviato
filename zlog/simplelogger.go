// +build !solution

// Leave an empty line above this comment.
package zlog

import (
	"fmt"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
	//. "github.com/uis-dat320-fall18/assignments/lab6" REMOVE
)

type Zaps []chzap.ChZap

func NewSimpleZapLogger() ZapLogger {
	zs := make(Zaps, 0)
	return &zs
}

func (zs *Zaps) LogZap(z chzap.ChZap) {
	*zs = append(*zs, z)
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
	// TODO uncomment this code when ToChan and FromChan added to ChZap struct
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

// Channels creates a slice of the channels found in the zaps(both to and from).
// ChZap format: {"2010/12/22, 20:22:32, 10.213.223.232, NRK2, NRK1", "20:22:32"}
func (zs *Zaps) Channels() []string {
	defer util.TimeElapsed(time.Now(), "simple.Channels")
	channels := make([]string, len(*zs))
	for _, channel := range *zs {
		channels = append(channels, channel.ToChan)
		channels = append(channels, channel.FromChan)
	}
	return channels
}

// ChannelsViewers creates a slice of ChannelViewers, which is defined in zaplogger.go.
// This is the number of viewers for each channel.
func (zs *Zaps) ChannelsViewers() []*ChannelViewers {
	defer util.TimeElapsed(time.Now(), "simple.ChannelsViewers")
	//TODO write this method (5p)
	// 	type ChannelViewers struct {
	// 		Channel string
	// 		Viewers int
	// }
	return nil
}
