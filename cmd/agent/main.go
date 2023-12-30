package main

import (
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
	"go.uber.org/zap"
)

const buffersNumber = 5

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(loggerZap *zap.Logger) {
		_ = loggerZap.Sync()
	}(logger)
	zap.ReplaceGlobals(logger)
	zap.S().Infoln(config.PrepBuildValues(BuildVersion, BuildDate, BuildCommit))
	sigs := initSigChan()

	cfg, err := config.LoadAgentConfig()
	if err != nil {
		zap.S().Fatalf("couldn't create a config %s", err)
	}
	logCfg(*cfg)
	initIfNeedAES(*cfg)

	if err = grpc.Init(*cfg); err != nil {
		zap.S().Infof("no gRPC by: %s", err)
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
				go metricsCase(metrics)

			case err = <-errCh:
				zap.S().Fatalf("cannot collect metrics: %s", err)

			case <-tickerReport.C:
				reportCase(cfg)
			}
		}
	}()
	zap.S().Infoln("awaiting a signal or press Ctrl+C to finish this agent")
	<-done
	zap.S().Infoln("exiting")
}

func metricsCase(metrics []model.MetricUnit) {
	for _, s := range metrics {
		repository.AddMetric(s)
	}
	zap.S().Infoln("added metrics:", len(metrics))
}

func initIfNeedAES(cfg config.AgentConfig) {
	if cfg.CryptoKey != "" {
		if err := repository.InitAgentAesKeys(cfg); err != nil {
			zap.S().Fatalf("couldn't init aes key %s", err)
		}
	}
}

func initSigChan() *chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return &sigs
}

func logCfg(cfg config.AgentConfig) {
	zap.S().Infoln("cfg:\n", "address:", cfg.Address, "\nkey:", cfg.HashKey)
	zap.S().Infoln("reportInterval:", cfg.ReportInterval, "\npollInterval:", cfg.PollInterval)
}

func gatherCase(metricsCh chan<- []model.MetricUnit, errCh chan error) {
	go gather.GetMemoryMetricsVariant(metricsCh, errCh)
	go gather.GetMetricsVariant(metricsCh, errCh)
}

func worker(workerID int64, cfg *config.AgentConfig, jobs <-chan []model.MetricUnit) {
	for j := range jobs {
		zap.S().Infoln("worker:", workerID, "started task:", j)
		// a real job.
		sender.ParcelsSend(cfg, j)

		zap.S().Infoln("worker:", workerID, "done task:", j)
	}
}

func reportCase(cfg *config.AgentConfig) {
	metrics := repository.GetAllMetrics()
	zap.S().Infoln("all jobs:", metrics)

	workerNumber := cfg.RateLimit // Rate limit
	if workerNumber == 0 {        // without workers
		sender.ParcelsSend(cfg, metrics)

		return
	}
	numJobs := len(metrics)
	jobs := make(chan []model.MetricUnit, numJobs)

	for w := int64(1); w <= workerNumber; w++ {
		go worker(w, cfg, jobs)
	}
	for j := 1; j <= numJobs; j++ {
		number := j - 1
		job := metrics[number:j]
		jobs <- job
	}
	close(jobs)
}
