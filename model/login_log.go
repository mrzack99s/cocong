package model

import "time"

type LoginLog struct {
	BaseModel

	TransactionAt time.Time
	IPAddress     string
	SuccessLogin  bool

	ByUser string
}
