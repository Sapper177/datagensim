package main

import (
	"context"
	"flag"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"github.com/Sapper177/datagensim/internal/sim"
)

func parseargs(cfg *sim.Config) {

	// Get Bus Name
	flag.StringVar(&cfg.BusName, "b", "MainBus", "Bus Name")

	// Get Source IP
	flag.StringVar(&cfg.SrcHost, "sh", "127.0.0.1", "Source Hostname or IP address")

	// Get Destination IP
	flag.StringVar(&cfg.DestHost, "dh", "127.0.0.1", "Destination Hostname or IP address")

	// Get Source Port
	flag.IntVar(&cfg.SrcPort, "sp", 0, "Source port")

	// Get Destination Port
	flag.IntVar(&cfg.DestPort, "dp", 0, "Destination port")

	// Optional Arguments
	flag.StringVar(&cfg.LogFile, "l", "/var/tmp/log", "Log file path")
	flag.StringVar(&cfg.LogLevel, "ll", "info", "Log level (debug, info, warn, error)")

	flag.Parse()
}

func main(){
	log.SetOutput(os.Stdout)

	// Create config and load from command line args
	cfg := new(sim.Config)
	parseargs(&cfg)

	// Create a new logger
	logger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "gosim")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logger)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Starting simulation with config:", cfg)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// Start the simulation in a goroutine
	go func() {
		if err := sim.Sim(ctx, cfg); err != nil {
			log.Println("Error in simulation:", err)
		}
	}()

	// Wait for OS signal
	sig := <-sigChan
	log.Println("Received signal:", sig)
	// Cancel the context to stop the simulation
	cancel()
	log.Println("Simulation stopped gracefully")

	// Close the logger
	if err := logger.Close(); err != nil {
		log.Println("Error closing logger:", err)
	}
}
