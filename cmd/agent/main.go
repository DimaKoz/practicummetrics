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

const buffersNumber = 5

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
	metricsCh := make(chan *[]model.MetricUnit, buffersNumber)
	defer close(metricsCh)
	errCh := make(chan error)
	defer close(errCh)
	done := make(chan bool)
	defer close(done)

	go func() {
		for {
			select {
			case <-sigs:
				done <- true

				return

			case <-tickerGathering.C:
				gatherCase(metricsCh, errCh)

			case metrics := <-metricsCh:
				go metricsCase(metrics, infoLog)

			case err = <-errCh:
				infoLog.Fatalf("cannot collect metrics: %s", err)

			case <-tickerReport.C:
				reportCase(cfg, infoLog)
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

func worker(workerID int64, cfg *config.AgentConfig, infoLog *log.Logger, jobs <-chan []model.MetricUnit) {
	for j := range jobs {
		infoLog.Println("worker:", workerID, "started task:", j)
		// a real job.
		sender.ParcelsSend(cfg, j)

		infoLog.Println("worker:", workerID, "done task:", j)
	}
}

func reportCase(cfg *config.AgentConfig, infoLog *log.Logger) {
	metrics := repository.GetAllMetrics()
	infoLog.Println("all jobs:", metrics)

	workerNumber := cfg.RateLimit // Rate limit
	if workerNumber == 0 {        // without workers
		sender.ParcelsSend(cfg, metrics)

		return
	}
	numJobs := len(metrics)
	jobs := make(chan []model.MetricUnit, numJobs)

	for w := int64(1); w <= workerNumber; w++ {
		go worker(w, cfg, infoLog, jobs)
	}
	for j := 1; j <= numJobs; j++ {
		number := j - 1
		job := metrics[number:j]
		jobs <- job
	}
	close(jobs)
}
