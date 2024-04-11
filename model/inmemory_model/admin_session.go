package inmemory_model

import (
	"time"

	"github.com/mrzack99s/cocong/model"
)

type AdminSession struct {
	model.BaseModel

	User              string
	AccessToken       string
	AccessTokenExpire time.Time

	RefreshToken       string
	RefreshTokenExpire time.Time
}
