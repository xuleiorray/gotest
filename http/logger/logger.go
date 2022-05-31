package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var LOGGER = log.New()

func init() {
	LOGGER.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:  true,
		EnvironmentOverrideColors: true,
		FullTimestamp:true,
		DisableLevelTruncation:true,
	})

    LOGGER.SetLevel(log.InfoLevel)

    logfile, _ := os.OpenFile("./perftest.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
    LOGGER.SetOutput(logfile)
    LOGGER.SetOutput(os.Stdout)

}
