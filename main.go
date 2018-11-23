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
	//	"github.com/AndrewBurian/powermux"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Flags
	var debug, quiet bool
	var port uint

	flag.BoolVar(&debug, "debug", false, "Debug verbosity")
	flag.BoolVar(&quiet, "quiet", false, "Errors only")
	flag.UintVar(&port, "port", 8080, "Port to bind to")
	flag.Parse()

	if debug && quiet {
		log.Fatal("Can only set one of -quiet and -debug")
	}

	if debug {
		log.SetLevel(log.DebugLevel)
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

	// kick of PD polling service
	go PollUpdates(updateStream)

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
	}()

	// Run server
	log.WithField("addr", server.Addr).Info("Server Starting")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("Error with server")
	}

	log.Info("Server Shutting down")
}
