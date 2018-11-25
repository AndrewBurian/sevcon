package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AndrewBurian/eventsource"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
)

type ConditionMonitor struct {
	currentLevel int
	latestEvent  *eventsource.Event
	client       *pagerduty.Client
}

func SetupMonitor(token string) *ConditionMonitor {
	return &ConditionMonitor{
		currentLevel: 5,
		client:       pagerduty.NewClient(token),
	}
}

func (mon *ConditionMonitor) PollUpdates(stream *eventsource.Stream, rate time.Duration) {

	// Events are created of type "update" starting from ID 2
	evFact := eventsource.EventIDFactory{
		Next: 2,
		NewFact: &eventsource.EventTypeFactory{
			Type: "update",
		},
	}

	// ID 1 is a fake event so we always have one for new clients
	mon.latestEvent = eventsource.DataEvent("5").ID("1")

	opts := pagerduty.ListIncidentsOptions{
		Statuses:  []string{"triggered", "acknowledged"},
		Urgencies: []string{"high"},
		Includes:  []string{"priorities"},
	}

	ticks := time.Tick(rate)
	for range ticks {

		log.Debug("Getting Incidents")
		response, err := mon.client.ListIncidents(opts)
		if err != nil {
			log.WithError(err).Error("Error querying PagerDuty")
			continue
		}

		var highest, cur int
		highest = 5
		log.WithField("count", len(response.Incidents)).Debug("Processing Incidents")
		for _, incident := range response.Incidents {
			if incident.Priority.ID == "" {
				log.Debug("Incident with no priority skipped")
				continue
			}
			n, err := fmt.Sscanf(incident.Priority.Summary, "SEV-%d", &cur)
			if err != nil || n != 1 {
				log.WithError(err).Error("Could not parse sev score")
				continue
			}

			if cur < highest {
				highest = cur
			}
		}

		if highest == mon.currentLevel {
			log.WithField("level", highest).Debug("No change")
			continue
		}

		mon.currentLevel = highest

		log.WithField("level", mon.currentLevel).Debug("Sending Update")
		ev := evFact.New()
		fmt.Fprintf(ev, "%d", mon.currentLevel)
		stream.Broadcast(ev)
		mon.latestEvent = ev
	}
}

func (mon *ConditionMonitor) NewClient(_ *http.Request, c *eventsource.Client) {
	log.Debug("New Client")
	if mon.latestEvent != nil {
		c.Send(mon.latestEvent)
	}
}

func DanceUpdates(stream *eventsource.Stream, rate time.Duration) {
	evFact := eventsource.EventIDFactory{
		Next: 1,
		NewFact: &eventsource.EventTypeFactory{
			Type: "update",
		},
	}
	ticks := time.Tick(rate)
	count := uint8(1)
	for range ticks {
		ev := evFact.New()
		fmt.Fprintf(ev, "%d", (count%5)+1)
		count++
		stream.Broadcast(ev)
	}
}
