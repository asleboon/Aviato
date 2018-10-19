// +build !solution

// Leave an empty line above this comment.
package zlog

import (
	"fmt"
	"time"
	. "github.com/uis-dat320/glabs/lab7"
)

type Zaps []ChZap

func NewSimpleZapLogger() ZapLogger {
	zs := make(Zaps, 0)
	return &zs
}

func (zs *Zaps) LogZap(z ChZap) {
	*zs = append(*zs, z)
}

func (zs *Zaps) Entries() int {
	return len(*zs)
}

func (zs *Zaps) String() string {
	return fmt.Sprintf("SS: %d", len(*zs))
}

// Viewers() returns the current number of viewers for a channel.
func (zs *Zaps) Viewers(chName string) int {
	defer TimeElapsed(time.Now(), "simple.Viewers")
	viewers := 0
	//TODO uncomment this code when ToChan and FromChan added to ChZap struct
	// for _, v := range *zs {
	// 	if v.ToChan == chName {
	// 		viewers++
	// 	}
	// 	if v.FromChan == chName {
	// 		viewers--
	// 	}
	// }
	return viewers
}

// Channels() creates a slice of the channels found in the zaps(both to and from).
func (zs *Zaps) Channels() []string {
	defer TimeElapsed(time.Now(), "simple.Channels")
	//TODO write this method (5p)
	return nil
}

// ChannelsViewers() creates a slice of ChannelViewers, which is defined in zaplogger.go.
// This is the number of viewers for each channel.
func (zs *Zaps) ChannelsViewers() []*ChannelViewers {
	defer TimeElapsed(time.Now(), "simple.ChannelsViewers")
	//TODO write this method (5p)
	return nil
}
