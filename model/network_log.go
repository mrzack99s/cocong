package model

import "time"

type NetworkLog struct {
	ID string

	TransactionAt time.Time
	Protocol      string

	SourceNetwork string
	SourcePort    string

	DestinationNetwork string
	DestinationPort    string

	TrafficFromInternet bool

	User string
}
