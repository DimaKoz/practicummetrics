package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/agent/send"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	flag2 "github.com/spf13/pflag"
	"runtime"
	"strconv"
	"time"
)

var pollInterval = time.Duration(2)
var reportInterval = time.Duration(10)
var alive time.Duration = 60
var Address string
var pFlag, rFlag string

func main() {

	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true

	flag2.StringVarP(&Address, "a", "a", "localhost:8080",
		"localhost:8080 by default")

	flag2.StringVarP(&pFlag, "p", "p", "2",
		"2 by default")

	flag2.StringVarP(&rFlag, "r", "r", "10",
		"10 by default")

	flag2.Parse()
	if s, err := strconv.ParseInt(pFlag, 10, 64); err == nil {
		pollInterval = time.Duration(s)
	} else {
		fmt.Println("pFlag:", pFlag, ", s:", s, ", err:", err)
	}

	if s, err := strconv.ParseInt(rFlag, 10, 64); err == nil {
		reportInterval = time.Duration(s)
	} else {
		fmt.Println("rFlag:", rFlag, ", s:", s, ", err:", err)
	}

	send.Address = Address

	fmt.Println("reportInterval:", reportInterval)
	fmt.Println("pollInterval:", pollInterval)

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
