package api_operation

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/services"
)

func (ctl *controller) loginLogQuery(c *gin.Context) {

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

	response := []model.LoginLog{}
	count, err := services.DBQuery(&response, offset, limit, search, or == "true", false)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}

func (ctl *controller) loginLogDump(c *gin.Context) {

	search := c.Query("search")
	or := c.Query("or")

	response := []model.LoginLog{}
	count, err := services.DBQuery(&response, 0, 0, search, or == "true", true)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}

func (ctl *controller) networkLogQuery(c *gin.Context) {

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

	response := []model.NetworkLog{}
	count, err := services.DBQuery(&response, offset, limit, search, or == "true", false)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}

func (ctl *controller) networkLogDump(c *gin.Context) {

	search := c.Query("search")
	or := c.Query("or")

	response := []model.NetworkLog{}
	count, err := services.DBQuery(&response, 0, 0, search, or == "true", true)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}

func (ctl *controller) logoutLogQuery(c *gin.Context) {

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

	response := []model.LogoutLog{}
	count, err := services.DBQuery(&response, offset, limit, search, or == "true", false)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}

func (ctl *controller) logoutLogDump(c *gin.Context) {

	search := c.Query("search")
	or := c.Query("or")

	response := []model.LogoutLog{}
	count, err := services.DBQuery(&response, 0, 0, search, or == "true", true)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"Count": count,
		"Data":  response,
	})

}
