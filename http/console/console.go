package console

import (
	"perftest/http/constants"
	"time"
)

type Console interface {
	Start()
	Stop(status constants.ConsoleStatus)
	IsReady()
	IsRunning()
	IsFinish()
	WaitTestFinish()
	Status()
	PrintTestReport()
}

type ConsoleRunner struct {
	Rendezvous    bool
	ConsoleStatus constants.ConsoleStatus
}

func NewConsoleRunner(rendezvous bool) (console *ConsoleRunner) {
	return &ConsoleRunner{
		Rendezvous:    rendezvous,
		ConsoleStatus: constants.Ready,
	}
}

func (console *ConsoleRunner) Start() {
	console.doStart()
	if console.Rendezvous {
		Context.simulcast()
	}
	console.changeStatus(constants.Running)

	// start test report stats with async mode
	Context.Stats.Start()
	
	// report at runtime
	Context.Stats.Export(Context.Reporter)
}

func (console *ConsoleRunner) doStart() {
	// create workers
	runCount := Context.RunCount
	workerCount := Context.WorkerCount
	avgRunCount := runCount / workerCount
	modRunCount := runCount % workerCount
	for i := 0; i < workerCount; i++ {
		if runCount > 0 {
			if avgRunCount > 0 {
				currentRunCount := avgRunCount
				if modRunCount > 0 {
					currentRunCount = avgRunCount + 1
					modRunCount--
				}
				Context.WorkerMgrs[i] = newWorkerMgr(i, console.Rendezvous, currentRunCount)
			} else { // only need run once for every worker
				Context.WorkerMgrs[i] = newWorkerMgr(i, console.Rendezvous, 1)
				runCount--
				if runCount == 0 {
					break
				}
			}
		} else {
			Context.WorkerMgrs[i] = newWorkerMgr(i, console.Rendezvous, -1)
		}
	}

	for _, wm := range Context.WorkerMgrs {
		wm.Start()
	}

}

func (console *ConsoleRunner) Stop(status constants.ConsoleStatus) {
	for _, wm := range Context.WorkerMgrs {
		wm.Stop()
	}
	console.changeStatus(status)
	Context.Stats.Stop()
}

func (console *ConsoleRunner) IsReady() bool {
	return console.ConsoleStatus == constants.Ready
}

func (console *ConsoleRunner) IsRunning() bool {
	return console.ConsoleStatus == constants.Running
}

func (console *ConsoleRunner) IsFinish() bool {
	return console.ConsoleStatus == constants.Finish || console.ConsoleStatus == constants.Stopped
}

func (console *ConsoleRunner) WaitTestFinish() {
	// wait for stress test complete
	if Context.TestType == constants.StressTest {
		<-time.After(Context.Duration)
	} else {
		Context.waitWorkers()
	}
	console.Stop(constants.Finish)
}

func (console *ConsoleRunner) Status() constants.ConsoleStatus {
	return console.ConsoleStatus
}

func (console *ConsoleRunner) changeStatus(status constants.ConsoleStatus) {
	console.ConsoleStatus = status
}

func (console *ConsoleRunner) ExportTestReport() {
	Context.Reporter.Await()
}
