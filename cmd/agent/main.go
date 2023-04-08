package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/agent/gather"
	"github.com/DimaKoz/practicummetrics/internal/agent/send"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	flag2 "github.com/spf13/pflag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var pollInterval = time.Duration(2)
var reportInterval = time.Duration(10)

const alive time.Duration = 20

var address string
var pFlag, rFlag string

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

	flagsProcess()

	tickerGathering := time.NewTicker(pollInterval * time.Second)
	done := make(chan bool)
	tickerReport := time.NewTicker(reportInterval * time.Second)
	go func() {
		for {
			select {

			case <-sigs:
				fmt.Println("sigs")
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

func flagsProcess() {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true

	flag2.StringVarP(&address, "a", "a", "localhost:8080",
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

	send.Address = address
	fmt.Println("reportInterval:", reportInterval)
	fmt.Println("pollInterval:", pollInterval)

}
