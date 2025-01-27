package model

import "time"

type LoginLog struct {
	BaseModel

	TransactionAt time.Time
	IPAddress     string
	Success       bool

	User string
}
