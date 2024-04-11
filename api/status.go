package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) getStatus(c *gin.Context) {

	clientIp := c.ClientIP()

	mSession := inmemory_model.Session{}
	err := vars.InMemoryDatabase.Where("ip_address = ?", clientIp).First(&mSession).Error
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if mSession.BandwidthID != nil {
		vars.Database.Where("id = ?", *mSession.BandwidthID).First(&mSession.Bandwidth)
	}

	var concurrent int64
	vars.InMemoryDatabase.Model(&inmemory_model.Session{}).Select("count(id)").Where("user = ?", mSession.User).Scan(&concurrent)

	c.HTML(200, "status.html", map[string]any{
		"ChangePasswordPageEndpoint": ChangePasswordPageEndpoint,
		"StatusEndpoint":             StatusEndpoint,
		"ErrorEndpoint":              ErrorEndpoint,
		"LogoutEndpoint":             LogoutEndpoint,
		"LogoutAllDeviceEndpoint":    LogoutAllDeviceEndpoint,
		"Concurrent":                 concurrent,
		"Session":                    mSession,
		"SessionLastSeen":            mSession.LastSeen.Format(time.RFC822),
	})

}
