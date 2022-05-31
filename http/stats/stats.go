package stats

import (
	"fmt"
	"perftest/http/model"
	"perftest/http/utils"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Stat struct {
	transactionChan     chan *model.Transaction // from
	transactions        []*model.Transaction    // to1
	reportChan 			chan *model.Report		// to2
}
var tickerStopFlag bool
var StatIns *Stat
var reportMutex sync.Mutex

func init() {
	StatIns = New()
}

func New() *Stat {
	if StatIns != nil {
		return StatIns
	}
	return &Stat{
		transactionChan:    make(chan *model.Transaction, 2048),
		transactions:       make([]*model.Transaction, 0, 1024),
		reportChan: 		make(chan *model.Report, 512),
	}
}

func (stat *Stat) Report(transaction *model.Transaction) {
	log.Info("Report one piece of transaction result.")
	stat.transactionChan <- transaction
}

func (stat *Stat) transfer() {
	for v := range stat.transactionChan {
		//log.Info("Transfer transaction to stat.transactions from stat.transactionChan.")
		reportMutex.Lock()
		stat.transactions = append(stat.transactions, v)
		log.Infof("stat.transactions length: %d", len(stat.transactions))
		reportMutex.Unlock()
	}
}

func (stat *Stat) Start() {
	log.Infoln("Start to collect performance test results.")
	go stat.transfer()

	go func(stat *Stat) {
		log.Info("Stats aggregate process startup.")
		var ticker = time.NewTicker(1 * time.Second)
		for t := range ticker.C {
			if tickerStopFlag {
				break
			}
			if stat.isEmpty() {
				log.Infoln("Stats result buffer is empty.")
				continue
			}
			log.Infof("Stats start to aggregate test results at %s", utils.FormatTime(t))
			stat.aggregate()
			log.Infof("Stats finish to aggregate test results at %s", utils.FormatTime(t))
		}
		log.Info("Stats aggregate process exit.")
	}(stat)
}

func (stat *Stat) Stop() {
	log.Infoln("Stats is about to stop in 5 secs")
	<-time.After(5*time.Second)
	tickerStopFlag = true
	close(stat.reportChan) // close report channel
	log.Infoln("Stats is stopped")
}

func (stat *Stat) aggregate() {
	reportMutex.Lock()
	transSum := len(stat.transactions)
	tempTrans := stat.transactions[:transSum]
	stat.transactions = stat.transactions[transSum:]
	reportMutex.Unlock()

	successCount, failedCount, totalTimeSpent := statusStats(tempTrans)
	tp50, tp90, tp99 := tpStats(tempTrans)
	avgRtt, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(totalTimeSpent)/float64(transSum)), 64)
	report := &model.Report{
		Period:       tempTrans[transSum-1].ExecEndTime.Sub(tempTrans[0].ExecStartTime),
		Total:        transSum,
		SuccessCount: successCount,
		FailedCount:  failedCount,

		AvgRTT: avgRtt,
		TP50:   tp50,
		TP90:   tp90,
		TP99:   tp99,
	}
	log.Infof("Stats aggregate multiple transations within one seconds into report: %s", utils.ToJSON(report))
	stat.reportChan <- report
}

func tpStats(trans []*model.Transaction) (tp50 int64, tp90 int64, tp99 int64) {
	sort.Slice(trans, func(i, j int) bool {
		return trans[i].RTT < trans[j].RTT
	})
	length := len(trans)
	
	tp50 = trans[length/2].RTT
	tp90 = trans[length * 9 / 10].RTT
	tp99 = trans[length * 99 / 100].RTT
	return
}

func statusStats(trans []*model.Transaction) (successCount int, failedCount int, totalTimeSpent int64) {
	for _, transaction := range trans {
		totalTimeSpent += transaction.RTT
		if transaction.Status {
			successCount++
		} else {
			failedCount++
		}
	}
	return
}

func (stat *Stat) isEmpty() bool {
	return len(stat.transactions) == 0
}

func (stat *Stat) Export(reporter *Reporter) {
	go reporter.output(stat.reportChan)
}