package zlog

import (
	"fmt"
	. "github.com/uis-dat320/glabs/lab7"
)

type ZapLogger interface {
	LogZap(z ChZap)
	Entries() int
	Viewers(channelName string) int
	Channels() []string
	ChannelsViewers() []*ChannelViewers
}

type ChannelViewers struct {
	Channel string
	Viewers int
}

func (cv ChannelViewers) String() string {
	return fmt.Sprintf("%s: %d", cv.Channel, cv.Viewers)
}

type ChanViewersList []*ChannelViewers
type ByViewers struct{ ChanViewersList }

func (t ChanViewersList) Len() int      { return len(t) }
func (t ChanViewersList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (s ByViewers) Less(i, j int) bool {
	return s.ChanViewersList[i].Viewers < s.ChanViewersList[j].Viewers
}
