package main

import (
	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/agent/sender"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("agent: couldn't create a config %s", err)
	}

	// from cfg:
	log.Println("agent: cfg:")
	log.Println("agent: address:", cfg.Address)
	log.Println("agent: reportInterval:", cfg.ReportInterval)
	log.Println("agent: pollInterval:", cfg.PollInterval)

	tickerGathering := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer tickerGathering.Stop()

	tickerReport := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer tickerReport.Stop()

	done := make(chan bool)
	go func() {
		for {
			select {

			case <-sigs:
				done <- true
				return

			case <-tickerGathering.C:
				metrics, err := gather.GetMetrics()
				if err != nil {
					log.Fatalf("agent: cannot collect metrics: %s", err)
				}
				for _, s := range *metrics {
					repository.AddMetric(s)
				}

			case <-tickerReport.C:
				metrics := repository.GetAllMetrics()
				sender.ParcelsSend(cfg, metrics)
			}
		}
	}()

	log.Println("agent: awaiting a signal or press Ctrl+C to finish this agent")

	<-done

	log.Println("agent: exiting")

}
