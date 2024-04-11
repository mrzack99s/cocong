package model

import "time"

type AdministratorLoginLog struct {
	BaseModel

	TransactionAt time.Time
	IPAddress     string
	SuccessLogin  bool

	ByUser string
}
