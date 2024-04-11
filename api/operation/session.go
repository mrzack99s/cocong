package api_operation

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/services"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) sessionQuery(c *gin.Context) {

	search := c.Query("search")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	or := c.Query("or")

	offset, e := strconv.Atoi(offsetStr)
	if e != nil {
		c.String(400, "offset is not correct, allow only integer")
		return
	}

	limit, e := strconv.Atoi(limitStr)
	if e != nil {
		c.String(400, "limit is not correct, allow only integer")
		return
	}

	response := []inmemory_model.Session{}
	count, err := services.DBQueryCustomDB(vars.InMemoryDatabase, &response, offset, limit, search, or == "true", false)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	for i, r := range response {
		if r.BandwidthID != nil {
			vars.Database.Where("id = ?", r.BandwidthID).First(&response[i].Bandwidth)
		}

	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}

func (ctl *controller) sessionKick(c *gin.Context) {

	var params struct {
		SessionID string
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	s := inmemory_model.Session{}
	if err := vars.InMemoryDatabase.Where("id = ?", params.SessionID).First(&s).Error; err != nil {
		c.String(500, "not found session")
		return
	}

	if err := session.CutOffSession(s); err != nil {
		c.String(500, "cannot cutoff session")
		return
	}

	c.Data(204, "text/plain", nil)
}
