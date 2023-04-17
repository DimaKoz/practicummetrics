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

	cfg, err := config.CreateConfig(config.ServerCfg)
	if err != nil {
		log.Fatalf("couldn't create a config %s", err)
	}

	// from cfg:
	log.Println("cfg:")
	log.Println("address:", cfg.Address)
	log.Println("reportInterval:", cfg.ReportInterval)
	log.Println("pollInterval:", cfg.PollInterval)

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
				metrics, err1 := gather.GetMetrics()
				if err1 != nil {
					log.Fatalf("cannot collect metrics: %s", err1)
				}
				for _, s := range metrics {
					repository.AddMetric(s)
				}

			case <-tickerReport.C:
				metrics := repository.GetAllMetrics()
				sender.ParcelsSend(cfg, metrics)
			}
		}
	}()

	log.Println("awaiting a signal or press Ctrl+C to finish this agent")

	<-done

	log.Println("exiting")

}
