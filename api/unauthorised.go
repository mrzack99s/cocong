package api

import (
	"github.com/gin-gonic/gin"
)

func (ctl *controller) unauthorised(c *gin.Context) {

	c.HTML(200, "unauthorised.html", nil)

}
