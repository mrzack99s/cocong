package types

import (
	"time"

	"github.com/mrzack99s/cocong/model"
)

type SessionInfo struct {
	ID string

	IPAddress string
	AuthType  string
	User      string
	LastSeen  time.Time

	BandwidthID *string
	Bandwidth   model.Bandwidth
}
