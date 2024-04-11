package inmemory_model

import (
	"time"

	"github.com/mrzack99s/cocong/model"
)

type Session struct {
	model.BaseModel

	IPAddress string
	AuthType  string
	User      string
	LastSeen  time.Time

	BandwidthID *string         `gorm:"type:varchar(36);"`
	Bandwidth   model.Bandwidth `gorm:"-"`
}
