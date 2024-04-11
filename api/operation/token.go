package api_operation

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) refreshToken(c *gin.Context) {

	var params struct {
		RefreshToken string
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	data, ok := vars.AdminSession.Get(params.RefreshToken)
	if !ok {
		c.String(400, "bad request!, invalid refresh token")
		return
	}

	dataStruct := data.(types.TokenSession)
	vars.AdminSession.Del(dataStruct.AccessToken)
	vars.AdminSession.Del(dataStruct.RefreshToken)

	accessToken := utils.SecretGenerator(128)
	refreshToken := utils.SecretGenerator(128)
	now := time.Now().In(vars.TZ)

	tokenSession := types.TokenSession{
		UserID:              dataStruct.UserID,
		AccessToken:         accessToken,
		AccessTokenExpired:  now.Add(time.Minute * 45),
		RefreshToken:        refreshToken,
		RefreshTokenExpired: now.Add(dataStruct.RefreshTokenExpired.Sub(now)),
	}

	vars.AdminSession.SetWithTTL(accessToken, tokenSession.UserID, 1, time.Minute*45)
	vars.AdminSession.SetWithTTL(refreshToken, tokenSession, 1, dataStruct.RefreshTokenExpired.Sub(now))

	c.JSON(200, tokenSession)

}

func (ctl *controller) me(c *gin.Context) {

	accessToken := c.Request.Header.Get("api-token")

	data, ok := vars.AdminSession.Get(accessToken)
	if !ok {
		c.String(500, "not found api-token")
		return
	}

	adm := model.Administrator{}
	if err := vars.Database.Where("id = ?", data).Omit("hashed").First(&adm).Error; err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, adm)

}

func (ctl *controller) logout(c *gin.Context) {

	accessToken := c.Request.Header.Get("api-token")
	vars.AdminSession.Del(accessToken)

	c.String(200, "ok")

}
