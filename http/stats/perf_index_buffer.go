package stats

import (
	"sort"
	"sync"
	"sync/atomic"
)

type PerfIndex struct {
	Total        uint32
	SuccessCount uint32
	FailedCount  uint32

	MaxRTT uint32
	MinRTT uint32
	AvgRTT float64
	TP50   uint32
	TP90   uint32
	TP99   uint32
}

var length uint32 = 256
var exponent uint8 = 8
var capacity uint32 = length * length
var maxIndex uint32 = capacity - 1

var perfIndexBuffer = &PerfIndexBuffer {
}

type PerfIndexBuffer struct {
	successCount uint32
	failedCount  uint32

	maxRTT	uint32
	minRTT	uint32

	timeCountMatrix [256][256]uint32
	timeCountLeakMap sync.Map
	leakMapKeyCounter uint32
}

func (buffer *PerfIndexBuffer) Success(rtt uint32, count uint32) {
	if count <= 0 {
		return
	}
	
	if rtt > buffer.maxRTT {
		buffer.maxRTT = rtt
	}

	if rtt < buffer.minRTT {
		buffer.minRTT = rtt
	}

	atomic.AddUint32(&buffer.successCount, 1)
	if rtt >= capacity {
		if v, err := buffer.timeCountLeakMap.Load(rtt); !err {
			buffer.timeCountLeakMap.Store(rtt, v.(uint32) + count)
		} else {
			buffer.timeCountLeakMap.Store(rtt, count)
			buffer.leakMapKeyCounter++
		}
		return
	}

	i := rtt >> exponent
	j := rtt & maxIndex
	atomic.AddUint32(&buffer.timeCountMatrix[i][j], count)
}

func (buffer *PerfIndexBuffer) Fail(count uint32) {
	if count <= 0 {
		return
	}
	atomic.AddUint32(&buffer.failedCount, count)
}

func (buffer *PerfIndexBuffer) Aggregate() (perfIndex PerfIndex) {
	perfIndex.Total = buffer.successCount + buffer.failedCount
	perfIndex.SuccessCount = buffer.successCount
	perfIndex.FailedCount = buffer.failedCount
	perfIndex.MaxRTT = buffer.maxRTT
	perfIndex.MinRTT = buffer.minRTT
	
	tp50Index, tp90Index, tp99Index := buffer.getTPXPos()

	var currPos uint32
	var prevPos uint32
	for i := 0; i < int(length); i++ {
		for j := 0; j < int(length); j++ {
			count := buffer.timeCountMatrix[i][j]
			if count > 0 {
				currPos = prevPos + count
				
				time := i * int(length) + j
				if tp50Index > prevPos && tp50Index <= currPos {
					perfIndex.TP50 = uint32(time)
				}
				if tp90Index > prevPos && tp90Index <= currPos {
					perfIndex.TP90 = uint32(time)
				}
				if tp99Index > prevPos && tp99Index <= currPos {
					perfIndex.TP99 = uint32(time)
				}
				prevPos = currPos
			}
		}
	}
	if tp50Index != 0 && tp90Index != 0 && tp99Index != 0 {
		return
	}
	
	for time := range buffer.sortMapKeys() {
		count, _ := buffer.timeCountLeakMap.Load(time)
		currPos = prevPos + count.(uint32)
		if tp50Index > prevPos && tp50Index <= currPos {
			perfIndex.TP50 = uint32(time)
		}
		if tp90Index > prevPos && tp90Index <= currPos {
			perfIndex.TP90 = uint32(time)
		}
		if tp99Index > prevPos && tp99Index <= currPos {
			perfIndex.TP99 = uint32(time)
		}
		if tp50Index != 0 && tp90Index != 0 && tp99Index != 0 {
			break
		}
		prevPos = currPos
	}

	return
}

func (buffer *PerfIndexBuffer) sortMapKeys() []int {
	keys := make([]int, 0, buffer.leakMapKeyCounter)
	buffer.timeCountLeakMap.Range(func(time, value any) bool {
		keys = append(keys, time.(int))
		return true
	})
	sort.Ints(keys)
	return keys
} 

func (buffer *PerfIndexBuffer) getTPXPos() (tp50 uint32, tp90 uint32, tp99 uint32) {
	tp50 = buffer.successCount * 50 / 100
	tp90 = buffer.successCount * 90 / 100
	tp99 = buffer.successCount * 99 / 100
	return
} 