package console

import "sync"

// worker manager to controll worker do task
type WorkerMgr struct{
	workerIndex 	int
	rendezvous 		bool
	runCount 		int
	runCountMutex 	sync.Mutex
}

func newWorkerMgr(workerIndex int, rendezvous bool, runCount int) (workerMgr *WorkerMgr) {
	workerMgr = new(WorkerMgr)
	workerMgr.workerIndex = workerIndex
	workerMgr.rendezvous = rendezvous
	workerMgr.runCount = runCount

	workerMgr.runCountMutex = sync.Mutex{}

	Context.registerWorkerMgr(workerIndex, workerMgr)
	return
}

func (workerMgr *WorkerMgr) Start() {
	Context.notifyWorkerStart()
	go workerMgr.doStart()
}

func (workerMgr *WorkerMgr) doStart() {
	defer func() {
		Context.notifyWorkerDone()  // notify console when worker complete
		Context.unRegisterWorkerMgr(workerMgr.workerIndex)
	}()

	if workerMgr.rendezvous {
		Context.standby() // wait for start signal
	}
	for workerMgr.hasNextRun() {
		Context.TestTask.DoTask()
	}
}

func (workerMgr *WorkerMgr) hasNextRun() (flag bool) {
	if workerMgr.runCount < 0 {
		return true
	}
	workerMgr.decrease()
	return workerMgr.runCount >= 0
}

func (workerMgr *WorkerMgr) decrease() {
	workerMgr.runCountMutex.Lock()
	workerMgr.runCount--
	workerMgr.runCountMutex.Unlock()
}

func (workerMgr *WorkerMgr) Stop() {
	workerMgr.runCountMutex.Lock()
	workerMgr.runCount = 0
	workerMgr.runCountMutex.Unlock()
}

