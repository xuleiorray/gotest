package task

import (
	"perftest/http/logger"
	"perftest/http/model"
	"perftest/http/stats"
	"perftest/http/utils"
	"time"
)

var log = logger.LOGGER

type PerfTestTask struct {
	HttpRequest	*model.HttpRequest
}

func (testTask *PerfTestTask) DoTask() {
	defer utils.Trace("DoTask")()
	httpReq := testTask.HttpRequest

	transaction := &model.Transaction{
		Name : httpReq.TransId,
		ExecStartTime: time.Now(),
		HttpRequest: testTask.HttpRequest,
	}
	httpResp := axiosClient.Dispatch(httpReq)
	log.Infof("% 4s %s %d", httpReq.Method, httpReq.Url, httpResp.StatusCode)

	transaction.ExecEndTime = time.Now()
	transaction.RTT = transaction.ExecEndTime.Sub(transaction.ExecStartTime).Milliseconds()
	transaction.Status = httpResp != nil

	testTask.report(transaction)
}

// submit one transaction to Stats
func (testTask *PerfTestTask) report(transaction *model.Transaction) {
	go stats.StatIns.Report(transaction)
}