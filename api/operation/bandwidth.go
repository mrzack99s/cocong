package api_operation

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/services"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) bandwidthCreate(c *gin.Context) {

	var params model.Bandwidth
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	if err := vars.Database.Where("name = ?", params.Name).First(&params).Error; err == nil {
		c.String(500, "duplicated profile")
		return
	}

	err := vars.Database.Create(&params).Error
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(201, "text/plain", nil)

}

func (ctl *controller) bandwidthDelete(c *gin.Context) {

	id := c.Query("id")
	if id == "" {
		c.String(400, "id required.")
	}

	err := vars.Database.Where("id = ?", id).Delete(&model.Bandwidth{}).Error
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(204, "text/plain", nil)

}

func (ctl *controller) bandwidthQuery(c *gin.Context) {

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

	response := []model.Bandwidth{}
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
