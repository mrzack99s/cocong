package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) getStatus(c *gin.Context) {

	clientIp := c.ClientIP()

	mSession, err := session.Instance.GetByIP(clientIp)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if mSession.BandwidthID != nil {
		vars.Database.Where("id = ?", mSession.BandwidthID).First(&mSession.Bandwidth)
	}

	concurrent := 0
	listUsernameSession, err := session.Instance.GetByUsername(mSession.User)
	if err == nil {
		concurrent = len(listUsernameSession)
	}

	c.HTML(200, "status.html", map[string]any{
		"ChangePasswordPageEndpoint": ChangePasswordPageEndpoint,
		"StatusEndpoint":             StatusEndpoint,
		"ErrorEndpoint":              ErrorEndpoint,
		"LogoutEndpoint":             LogoutEndpoint,
		"LogoutAllDeviceEndpoint":    LogoutAllDeviceEndpoint,
		"Concurrent":                 concurrent,
		"Session":                    mSession,
		"SessionLastSeen":            mSession.LastSeen.Format(time.RFC822),
		"RedirectURL":                vars.Config.RedirectURL,
	})

}
