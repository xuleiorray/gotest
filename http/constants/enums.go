package constants

type TestType uint8
const(
	SpikeTest TestType = iota + 1
	StressTest
	LoadTest
)

type ConsoleStatus uint8
const(
	Ready ConsoleStatus = iota + 1
	Running
	Finish
	Stopped
)