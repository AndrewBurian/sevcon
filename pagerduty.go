package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AndrewBurian/eventsource"
	log "github.com/sirupsen/logrus"
)

type ConditionMonitor struct {
	currentLevel uint
	latestEvent  *eventsource.Event
}

func (mon *ConditionMonitor) PollUpdates(stream *eventsource.Stream) {

	evFact := eventsource.EventIDFactory{
		Next: 1,
	}

	ticks := time.Tick(time.Second)
	count := uint8(1)
	for _ = range ticks {
		ev := evFact.New()
		fmt.Fprintf(ev, "%d", (count%5)+1)
		count++
		stream.Broadcast(ev)
		mon.latestEvent = ev
	}
}

func (mon *ConditionMonitor) NewClient(_ *http.Request, c *eventsource.Client) {
	log.Debug("New Client")
	c.Send(mon.latestEvent)
}
