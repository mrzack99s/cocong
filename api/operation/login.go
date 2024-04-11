package api_operation

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) login(c *gin.Context) {

	var params types.CredentialVerification
	if err := c.ShouldBindJSON(&params); err != nil {
		c.String(400, "bad request!")
		return
	}

	admin := model.Administrator{}
	err := vars.Database.
		Where("username = ? and hashed = ?", params.Username, utils.Sha512encode(params.Password)).
		First(&admin).Error
	if err != nil {
		c.String(500, "credential incorrect")
		return
	}

	if !admin.Enable {
		c.String(500, "this account has been disable")
		return
	}

	accessToken := utils.SecretGenerator(128)
	refreshToken := utils.SecretGenerator(128)
	now := time.Now().In(vars.TZ)

	tokenSession := types.TokenSession{
		UserID:              admin.ID,
		AccessToken:         accessToken,
		AccessTokenExpired:  now.Add(time.Minute * 45),
		RefreshToken:        refreshToken,
		RefreshTokenExpired: now.Add(time.Hour * 24),
	}

	vars.AdminSession.SetWithTTL(accessToken, admin.ID, 1, time.Minute*45)
	vars.AdminSession.SetWithTTL(refreshToken, tokenSession, 1, time.Hour*24)

	c.JSON(200, tokenSession)

}
