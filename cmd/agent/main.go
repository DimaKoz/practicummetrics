package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/agent/send"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPollInterval   = time.Duration(2)
	defaultReportInterval = time.Duration(10)
	defaultAddress        = "localhost:8080"
	alive                 = time.Duration(21)
)

func sleepingKiller() {
	var stopTime = alive * time.Second
	time.Sleep(stopTime)
	fmt.Println("Time to die")
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg := &model.Config{}
	config.AgentInitConfig(cfg, defaultAddress, defaultReportInterval, defaultPollInterval)

	// from cfg:
	fmt.Println("cfg:")
	fmt.Println("address:", cfg.Address)
	fmt.Println("reportInterval:", cfg.ReportInterval)
	fmt.Println("pollInterval:", cfg.PollInterval)
	send.Address = cfg.Address
	tickerGathering := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	done := make(chan bool)
	go func() {
		for {
			select {

			case t := <-sigs:
				fmt.Println("sigs:", t)
				done <- true
				return

			case t := <-tickerGathering.C:
				fmt.Println("gathering info Tick at", t)
				metrics := gather.GetMetrics()
				for _, s := range *metrics {
					repository.AddMetricMemStorage(s)
				}

			case t := <-tickerReport.C:
				fmt.Println("sending info Tick at", t)
				metrics := repository.GetMetricsMemStorage()
				send.ParcelsSend(metrics)
			}
		}
	}()

	go sleepingKiller()

	fmt.Println("awaiting signal")
	<-done
	tickerGathering.Stop()
	tickerReport.Stop()

	fmt.Println("exiting")

}
