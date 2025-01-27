package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/utils"
)

func (ctl *controller) authLoginPage(c *gin.Context) {

	clientIp := c.ClientIP()

	_, err := session.Instance.GetByIP(clientIp)
	if err == nil {
		c.Redirect(http.StatusFound, "/status")
		return
	}

	image, copyright, _ := utils.FetchBingImage()

	c.HTML(200, "login.html", map[string]any{
		"StatusEndpoint": StatusEndpoint,
		"ErrorEndpoint":  ErrorEndpoint,
		"AuthEndpoint":   AuthEndpoint,
		"BGImage":        image,
		"BGCopyright":    copyright,
	})

}
