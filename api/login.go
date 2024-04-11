package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) authLoginPage(c *gin.Context) {

	clientIp := c.ClientIP()

	mSession := inmemory_model.Session{}
	err := vars.InMemoryDatabase.Where("ip_address = ?", clientIp).First(&mSession).Error
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
