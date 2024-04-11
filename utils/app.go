package utils

import (
	"os"
	"strings"
	"time"

	"github.com/mrzack99s/cocong/vars"
)

type AppCredentialsType struct {
	APIToken string `yaml:"api_token"`
}

var (
	AppCredentials AppCredentialsType = AppCredentialsType{}
)

func VerifyTimeZone() {

	if !ValidTimeZone(vars.Config.TimeZone) {
		vars.SystemLog.Println("timezone is not corrent")
		os.Exit(0)
	}

	vars.TZ, _ = time.LoadLocation(vars.Config.TimeZone)
}

func ValidTimeZone(tz string) bool {
	for _, v := range vars.AllowTZ {
		if strings.TrimSpace(v) == strings.TrimSpace(tz) {
			return true
		}
	}

	return false

}

func GetTimeZone() (loc *time.Location) {
	loc, _ = time.LoadLocation(vars.Config.TimeZone)
	return
}
