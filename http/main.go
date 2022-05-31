package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"perftest/http/console"
	"perftest/http/constants"
	"perftest/http/logger"
	"perftest/http/model"
	"perftest/http/stats"
	"perftest/http/task"
	"perftest/http/utils"
	"time"
)

var log = logger.LOGGER

const(
	START = iota + 1
	STOP
)

func main() {
	log.Infof("Welcome to use HTTP load test tool, %s", utils.FormatNow())

	var runCount 		int
	var duration 		time.Duration
	var workerCount		int
	var rendezvous 		bool
	var testConfig		string
	var reportPath		string
	// init variables above declared
	flag.IntVar(&runCount, "rc", -1, "run count")
	flag.DurationVar(&duration, "duration", time.Duration(0), "duration to run test")
	flag.IntVar(&workerCount, "wc", 1, "worker count")
	flag.BoolVar(&rendezvous, "sync", true, "do test concurrently")
	flag.StringVar(&testConfig, "config", "./gotest.json", "config file for testing")
	flag.StringVar(&reportPath, "reportPath", "./testReport.txt", "report filepath for testing")
	flag.Parse()

	var testType		constants.TestType
	if runCount > 0 {
		testType = constants.SpikeTest
	} else if duration.Seconds() > 0 {
		testType = constants.StressTest
	} else {
		interactOption(&runCount, &duration, &workerCount, &testType, &testConfig)
	}
	log.Infof("console params: RunCount:%d, Duration:%d, WorkerCount:%d, TestType:%d",
		runCount,
		duration,
		workerCount,
		testType)

	// init console context
	console.Context.RunCount = runCount
	console.Context.Duration = duration 
	console.Context.WorkerCount = workerCount
	console.Context.TestType = testType
	httpRequest, err := model.FromJsonFile("./gotest.json")
	if err != nil {
		log.Fatalf("Failed to read json file, %w", err)
	}
	console.Context.TestTask = task.PerfTestTask {
		HttpRequest: httpRequest,
	}

	// report file
	reportFile, err := os.OpenFile(reportPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		log.Fatalf("open report file error, %s", err.Error())
		return
	}
	console.Context.Reporter = &stats.Reporter{
		ReportFile: reportFile,
		FinishSignal: make(chan interface{}, 1),
	}

	// init console runner
	console := console.NewConsoleRunner(rendezvous)

	// launch test
	log.Infof("Start console runner...")
	console.Start()
	log.Infof("console runner...")

	// wait here until all worker finish their work
	console.WaitTestFinish()
	
	console.ExportTestReport()
}

func interactOption(runCount *int, duration *time.Duration, workerCount *int, testType *constants.TestType, config *string) {

	fmt.Println("Please provide some required parameters before run performance test.")
	
	fmt.Println("Select test type, 1: spike test, 2: stress test.")
	fmt.Scanln(testType)

	if *testType == 1 {
		fmt.Println("run count: ")
		fmt.Scanln(runCount)
	} else {
		fmt.Println("run duration, such as \"300ms\", \"15s\" or \"2h45m\".")
		var temp string
		fmt.Scanln(&temp)
		*duration,_  = time.ParseDuration(temp)
	}

	fmt.Println("Input worker count: ")
	fmt.Scanln(workerCount)

	fmt.Println("Input test config file, default value: './gotest.json'.")
	fmt.Scanln(config)
	file, err := os.Open(*config)
	if errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("Config file you provided could not be found, %s", *config)
	}
	defer file.Close()
}

