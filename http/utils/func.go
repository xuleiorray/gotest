package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(i interface{}, seps ...rune) string {
	// perftest/http/task.(*PerfTestTask).DoTask-fm
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		if size == 1 {
			return fields[0]
		}
		return fields[size-2]
	}
	return ""
}