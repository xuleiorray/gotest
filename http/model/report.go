package model

import "time"

type Transaction struct {
	Name      		string
	ExecStartTime 	time.Time
	ExecEndTime 	time.Time
	RTT				int64	// average of response time
	Status 			bool			// transaction status

	LogTime			int64 // time when this transaction record generated

	HttpRequest 	*HttpRequest
}

type Report struct {

	Period			int64

	Transactions	[]*Transaction
	Total			int
	SuccessCount 	int
	FailedCount		int
	
	AvgRTT			float64
	TP50			int64
	TP90			int64
	TP99			int64
	
}