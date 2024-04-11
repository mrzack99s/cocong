package api

import (
	"github.com/gin-gonic/gin"
)

func (ctl *controller) authErrorPage(c *gin.Context) {

	errMsg := c.Query("msg")
	c.HTML(200, "error.html", map[string]any{
		"ErrorMsg":       errMsg,
		"LoginEndpoint":  LoginEndpoint,
		"StatusEndpoint": StatusEndpoint,
		"ErrorEndpoint":  ErrorEndpoint,
	})

}
