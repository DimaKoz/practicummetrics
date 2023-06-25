package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/agent/sender"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
)

func main() {
	infoLog := log.Default()
	infoLog.SetPrefix("agent: INFO: ")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("couldn't create a config %s", err)
	}

	// from cfg:
	infoLog.Println("cfg:\n", "address:", cfg.Address, "\nkey:", cfg.HashKey)
	infoLog.Println("reportInterval:", cfg.ReportInterval)
	infoLog.Println("pollInterval:", cfg.PollInterval)

	tickerGathering := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer tickerGathering.Stop()

	tickerReport := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer tickerReport.Stop()
	metricsCh := make(chan *[]model.MetricUnit)
	errCh := make(chan error)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-sigs:
				done <- true

				return

			case <-tickerGathering.C:
				gatherCase(metricsCh, errCh)

			case metrics := <-metricsCh:
				metricsCase(metrics, infoLog)

			case err = <-errCh:
				infoLog.Fatalf("cannot collect metrics: %s", err)

			case <-tickerReport.C:
				reportCase(cfg)
			}
		}
	}()

	infoLog.Println("awaiting a signal or press Ctrl+C to finish this agent")

	<-done

	infoLog.Println("exiting")
}

func metricsCase(metrics *[]model.MetricUnit, infoLog *log.Logger) {
	for _, s := range *metrics {
		repository.AddMetric(s)
	}
	infoLog.Println("added metrics:", len(*metrics))
}

func gatherCase(metricsCh chan *[]model.MetricUnit, errCh chan error) {
	go gather.GetMemoryMetrics(metricsCh, errCh)
	go gather.GetMetrics(metricsCh, errCh)
}

func reportCase(cfg *config.AgentConfig) {
	metrics := repository.GetAllMetrics()
	sender.ParcelsSend(cfg, metrics)
}
