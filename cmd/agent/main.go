package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/agent/send"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"runtime"
	"time"
)

var pollInterval time.Duration = 2
var reportInterval time.Duration = 10
var alive time.Duration = 60

func main() {

	tickerGathering := time.NewTicker(pollInterval * time.Second)
	done := make(chan bool)
	tickerReport := time.NewTicker(reportInterval * time.Second)
	go func() {
		for {
			select {

			case <-done:
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

	var stopTime = alive * time.Second
	time.Sleep(stopTime)
	tickerGathering.Stop()
	tickerReport.Stop()
	done <- true
	fmt.Println("Ticker stopped")

}

func printMemStats(message string, rtm runtime.MemStats) {
	fmt.Println("\n===", message, "===")

	fmt.Println("Mallocs: ", rtm.Mallocs)

	fmt.Println("LiveObjects: ", rtm.Mallocs-rtm.Frees)
	fmt.Println("PauseTotalNs: ", rtm.PauseTotalNs)
	fmt.Println("NumGC: ", rtm.NumGC)
	fmt.Println("LastGC: ", time.UnixMilli(int64(rtm.LastGC/1_000_000)))
	fmt.Println("HeapObjects: ", rtm.HeapObjects)
	fmt.Println("HeapAlloc: ", rtm.HeapAlloc)
}
