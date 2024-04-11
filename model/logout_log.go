package model

import "time"

type LogoutLog struct {
	BaseModel

	TransactionAt time.Time
	IPAddress     string
	ByUser        string
}
