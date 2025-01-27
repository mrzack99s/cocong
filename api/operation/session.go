package api_operation

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) sessionQuery(c *gin.Context) {

	search := c.Query("search")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
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
	response, err := session.Instance.Search(search, offset, limit)
	if err != nil || response.Count == 0 {
		c.JSON(200, gin.H{
			"Count": 0,
			"Data":  []any{},
		})
		return
	}

	for i, r := range response.Data {
		if r.BandwidthID != nil {
			vars.Database.Where("id = ?", r.BandwidthID).First(&response.Data[i].Bandwidth)
		}

	}

	c.JSON(200, response)

}

func (ctl *controller) sessionKick(c *gin.Context) {

	var params struct {
		SessionID string
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	// _, err := session.Instance.GetByID(params.SessionID)
	// if err != nil {
	// 	c.String(500, "not found session")
	// 	return
	// }

	err := session.Instance.DeleteByID(params.SessionID)
	if err != nil {
		c.String(500, "not found session")
		return
	}

	c.Data(204, "text/plain", nil)
}
