package utils

import (
	"perftest/http/logger"
	"time"
)

var log = logger.LOGGER

/**
 * defer Trace("method_name")()
 */
func Trace(msg string) func() {
	start := time.Now()
    log.Infof("Enter (%s) time start at, %s", msg, FormatTime(start))
    return func() { 
		log.Infof("Exit (%s) time end at, %s, duration: %dms", msg, FormatTime(time.Now()), time.Since(start).Milliseconds())
    }
}

func TraceFunc(fun func()) func(...interface{}) {
	return func(...interface{}) {
		start := time.Now()
		log.Infof("Enter (%s) time start at, %s", GetFuncName(fun, '/', '.', '-'), FormatTime(start))
		fun()
		log.Infof("Exit (%s) time end at, %s, duration: %dms", GetFuncName(fun, '/', '.', '-'), FormatTime(time.Now()), time.Since(start).Milliseconds())
	}
}