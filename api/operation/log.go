package api_operation

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/services"
	"github.com/mrzack99s/cocong/vars"
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

	query := bleve.NewMatchQuery(search)
	searchRequest := bleve.NewSearchRequest(query)

	// กำหนด From และ Size สำหรับ Pagination
	searchRequest.From = offset
	searchRequest.Size = limit

	// ค้นหาข้อมูล
	searchResult, err := vars.NetLogDatabase.Search(searchRequest)
	if err != nil {
		c.String(500, fmt.Sprintf("Search error: %v", err))
		return
	}

	result_doc := []model.NetworkLog{}
	for _, hit := range searchResult.Hits {

		docBytes, err := vars.NetLogDatabase.GetInternal([]byte(hit.ID))
		if err != nil {
			c.String(500, "Error fetching document")
			return
		}

		if docBytes == nil {
			c.String(500, "Document not found")
			return
		}

		// แปลงข้อมูลจาก JSON เป็น Document struct
		var doc model.NetworkLog
		err = json.Unmarshal(docBytes, &doc)
		if err != nil {
			c.String(500, "Error unmarshalling document")
			return
		}

		result_doc = append(result_doc, doc)
	}

	// response := []model.NetworkLog{}
	// count, err := services.DBQueryCustomDB(vars.NetLogDatabaseRO, &response, offset, limit, search, or == "true", false)
	// if err != nil {
	// 	c.String(500, err.Error())
	// 	return
	// }

	c.JSON(200, gin.H{
		"Count": searchResult.Total,
		"Data":  result_doc,
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
