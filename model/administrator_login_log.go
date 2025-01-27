package model

import "time"

type AdministratorLoginLog struct {
	BaseModel

	TransactionAt time.Time
	IPAddress     string
	Success       bool

	User string
}
