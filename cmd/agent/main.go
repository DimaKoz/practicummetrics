package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/agent/grpc"
	"github.com/DimaKoz/practicummetrics/internal/agent/sender"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
)

const buffersNumber = 5

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	infoLog := initInfoLogger()

	sigs := initSigChan()

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		log.Fatalf("couldn't create a config %s", err)
	}
	logCfg(infoLog, *cfg)
	initIfNeedAES(*cfg)

	if err = grpc.Init(*cfg, infoLog); err != nil {
		infoLog.Printf("no gRPC by: %s", err)
	}
	defer grpc.Close()
	tickerGathering := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer tickerGathering.Stop()

	tickerReport := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer tickerReport.Stop()
	metricsCh := make(chan []model.MetricUnit, buffersNumber)
	defer close(metricsCh)
	errCh := make(chan error)
	defer close(errCh)
	done := make(chan bool)
	defer close(done)
	go func() {
		for {
			select {
			case <-*sigs:
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

func metricsCase(metrics []model.MetricUnit, infoLog *log.Logger) {
	for _, s := range metrics {
		repository.AddMetric(s)
	}
	infoLog.Println("added metrics:", len(metrics))
}

func initIfNeedAES(cfg config.AgentConfig) {
	if cfg.CryptoKey != "" {
		if err := repository.InitAgentAesKeys(cfg); err != nil {
			log.Fatalf("couldn't init aes key %s", err)
		}
	}
}

func initInfoLogger() *log.Logger {
	infoLog := log.Default()
	infoLog.SetPrefix("agent: INFO: ")
	infoLog.Println(config.PrepBuildValues(BuildVersion, BuildDate, BuildCommit))

	return infoLog
}

func initSigChan() *chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return &sigs
}

func logCfg(iLog *log.Logger, cfg config.AgentConfig) {
	// from cfg:
	iLog.Println("cfg:\n", "address:", cfg.Address, "\nkey:", cfg.HashKey)
	iLog.Println("reportInterval:", cfg.ReportInterval, "\npollInterval:", cfg.PollInterval)
}

func gatherCase(metricsCh chan<- []model.MetricUnit, errCh chan error) {
	go gather.GetMemoryMetricsVariant(metricsCh, errCh)
	go gather.GetMetricsVariant(metricsCh, errCh)
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
