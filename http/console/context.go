package console

import (
	"perftest/http/constants"
	"perftest/http/stats"
	"perftest/http/task"
	"sync"
	"time"
)

type ConsoleContext struct {
	TestTask 			task.PerfTestTask
	TestType 			constants.TestType

	WorkerRunCount 		map[int]int   		// only for SpikeTest
	RunCount			int					// only for SpikeTest
	Duration 			time.Duration 		// only for StressTest
	
	WorkerCount			int			
	WorkerMgrs 			map[int]*WorkerMgr

	WorkerCond			*sync.Cond 			// signal to start for all workers
	WorkerMonitor		*sync.WaitGroup 	// monitor worker start and done
	
	Stats				*stats.Stat
	Reporter			*stats.Reporter
}

var (
	Context *ConsoleContext
	standbyWorkers chan struct{} = make(chan struct{})
)
func init() {
	Context = new(ConsoleContext)
	Context.TestType = constants.SpikeTest
	Context.WorkerRunCount = make(map[int]int)
	Context.WorkerMgrs = make(map[int]*WorkerMgr)

	Context.WorkerCond = sync.NewCond(&sync.Mutex{})
	Context.WorkerMonitor = &sync.WaitGroup{}

	Context.Stats = stats.StatIns
}

func (context *ConsoleContext) registerWorkerMgr(index int, workerMgr *WorkerMgr) {
	context.WorkerMgrs[index] = workerMgr
}
func (context *ConsoleContext) unRegisterWorkerMgr(index int) {
	delete(context.WorkerMgrs, index)
}

func (context *ConsoleContext) standby() { // wait
	context.WorkerCond.L.Lock()
	standbyWorkers <- struct{}{}
	context.WorkerCond.Wait()
	context.WorkerCond.L.Unlock()
}
func (context *ConsoleContext) simulcast() { // release all workers 
	for i:=0; i < context.WorkerCount; i++ {
		<-standbyWorkers
	}
	context.WorkerCond.L.Lock()
	context.WorkerCond.Broadcast()
	context.WorkerCond.L.Unlock()
}

func (context *ConsoleContext) notifyWorkerStart() {
	context.WorkerMonitor.Add(1)
}
func (context *ConsoleContext) notifyWorkerDone() {
	context.WorkerMonitor.Done()
}
func (context *ConsoleContext) waitWorkers() {
	context.WorkerMonitor.Wait()
}