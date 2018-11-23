package main

import (
	"fmt"
	"time"

	"github.com/AndrewBurian/eventsource"
)

func PollUpdates(stream *eventsource.Stream) {

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
	}
}
