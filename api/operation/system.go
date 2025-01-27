package api_operation

import (
	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/services"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) getConfig(c *gin.Context) {
	c.JSON(200, vars.Config)
}

func (ctl *controller) writeConfig(c *gin.Context) {

	var params types.ConfigType
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	services.ConfigUpdate(params)
	services.WriteConfigUpdate()

	c.Data(204, "text/plain", nil)
}

func (ctl *controller) restartService(c *gin.Context) {
	services.RestartService()
}

func (ctl *controller) restartDNSServer(c *gin.Context) {
	services.RestartDNSServer()
}
