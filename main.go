package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndrewBurian/eventsource"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Flags
	var debug, quiet, test, help bool
	var port uint
	var rate string

	flag.BoolVar(&debug, "debug", false, "Debug verbosity")
	flag.BoolVar(&quiet, "quiet", false, "Errors only")
	flag.UintVar(&port, "port", 8080, "Port to bind to")
	flag.BoolVar(&test, "test", false, "Rotates status on a pattern instead of using real data")
	flag.BoolVar(&help, "help", false, "Display usage")
	flag.StringVar(&rate, "rate", "1m", "Update polling rate")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if debug && quiet {
		log.Fatal("Can only set one of -quiet and -debug")
	}

	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Running at debug verbosity")
	} else if quiet {
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	mux := http.NewServeMux()

	// Static site
	fileServer := http.FileServer(http.Dir("./site"))
	mux.Handle("/", fileServer)

	// SSE Update stream
	updateStream := eventsource.NewStream()
	mux.Handle("/updates", updateStream)

	// Polling rate
	pollRate, err := time.ParseDuration(rate)
	if err != nil {
		log.WithError(err).Fatal("Unable to parse polling rate")
	}

	if test {
		log.Warn("Running in test mode, no real data is being used")
		go DanceUpdates(updateStream, pollRate)
	} else {
		token, found := os.LookupEnv("PAGERDUTY_TOKEN")
		if !found {
			log.Fatal("Need PAGERDUTY_TOKEN env variable")
		}

		// Setup Monitor
		mon := SetupMonitor(token)

		// kick of PD polling service
		go mon.PollUpdates(updateStream, pollRate)
		updateStream.ClientConnectHook(mon.NewClient)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port),
		Handler: mux,
	}

	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		log.WithField("signal", s).Info("Trapped signal")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("Error shutting down server")
		}
		updateStream.Shutdown()
	}()

	// Run server
	log.WithField("addr", server.Addr).Info("Server Starting")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("Error with server")
	}

	log.Info("Server Shutting down")
}
