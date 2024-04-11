package types

import "time"

type NetLogType struct {
	Timestamp time.Time
	Proto     string
	SrcAddr   string
	DstAddr   string
	SrcPort   string
	DstPort   string
	NacID     string
}
