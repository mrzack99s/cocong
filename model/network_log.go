package model

import "time"

type NetworkLog struct {
	BaseModel

	TransactionAt time.Time
	Protocol      string

	SourceNetwork string
	SourcePort    string

	DestinationNetwork string
	DestinationPort    string

	TrafficFromInternet bool

	ByUser string
}
