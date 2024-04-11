package types

import "time"

type TokenSession struct {
	UserID              string
	AccessToken         string
	AccessTokenExpired  time.Time
	RefreshToken        string
	RefreshTokenExpired time.Time
}
