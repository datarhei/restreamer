package parse

import (
	"time"

	"github.com/prep/average"
)

type averager struct {
	fps     *average.SlidingWindow
	pps     *average.SlidingWindow
	bitrate *average.SlidingWindow
}

func (a *averager) init(window, granularity time.Duration) {
	a.fps, _ = average.New(window, granularity)
	a.pps, _ = average.New(window, granularity)
	a.bitrate, _ = average.New(window, granularity)
}

func (a *averager) stop() {
	a.fps.Stop()
	a.pps.Stop()
	a.bitrate.Stop()
}
