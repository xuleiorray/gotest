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
	ReportFile   *os.File
	FinishSignal chan interface{}
}

func (reporter *Reporter) output(reportChan chan *model.Report) {
	defer reporter.release()

	log.Infoln("Start to output report into local file.")
	bufWriter := bufio.NewWriter(reporter.ReportFile)

	for report := range reportChan {
		bufWriter.WriteString(utils.ToJSON(report) + "\n")
		bufWriter.Flush()
	}
	log.Infoln("Finish to output report into local file.")
}

func (reporter *Reporter) release() {
	reporter.printSummary()
	
	log.Infof("Close write report to file: %s", reporter.ReportFile.Name())
	reporter.ReportFile.Close()
	reporter.FinishSignal <- struct{}{}
}

// Wait for exporting test report to file complete and close normally
func (reporter *Reporter) Await() {
	log.Infof("Wait reporter output finish.")
	<-reporter.FinishSignal
	log.Infof("Reporter output finish.")
}

func (reporter *Reporter) printSummary() {
	log.Infof("^_^^_^^_^Performance Test Statistics^_^^_^^_^:")
	perfIndex := perfIndexBuffer.Aggregate()

	log.Infof(utils.ToJSON(perfIndex))
}
