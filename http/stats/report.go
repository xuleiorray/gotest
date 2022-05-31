package stats

import (
	"bufio"
	"os"
	"perftest/http/logger"
	"perftest/http/model"
	"perftest/http/utils"
)
var log = logger.LOGGER

type Reporter struct {
	ReportFile		*os.File
	FinishSignal  	chan interface{}
}

func (reporter *Reporter) output(reportChan chan *model.Report) {
	defer reporter.release()
	
	log.Infoln("Start to output report into local file.")
	bufWriter := bufio.NewWriter(reporter.ReportFile)
	
	for report := range reportChan {
		bufWriter.WriteString(utils.ToJSON(report)+"\n")
		bufWriter.Flush()
	}
	log.Infoln("Finish to output report into local file.")
}

func (reporter *Reporter) release() {
	log.Infof("Close write report to file: %s", reporter.ReportFile.Name())
	reporter.ReportFile.Close()
	reporter.FinishSignal <- struct{}{}
}
func (reporter *Reporter) Await() {
	log.Infof("Wait reporter output finish.")
	<-reporter.FinishSignal
	log.Infof("Reporter output finish.")
}