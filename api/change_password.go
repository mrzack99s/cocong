package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/session"
)

func (ctl *controller) changePasswordPage(c *gin.Context) {

	clientIp := c.ClientIP()

	mSession, err := session.Instance.GetByIP(clientIp)
	if err != nil {
		c.Redirect(http.StatusFound, "/error?msg=Please login first!")
		return
	}

	c.HTML(200, "chpassword.html", map[string]any{
		"LoginEndpoint":          LoginEndpoint,
		"StatusEndpoint":         StatusEndpoint,
		"ErrorEndpoint":          ErrorEndpoint,
		"ChangePasswordEndpoint": ChangePasswordEndpoint,
		"Username":               mSession.User,
	})

}
