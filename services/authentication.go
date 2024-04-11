package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func Authentication(ctx *gin.Context, cred types.CredentialVerification) (authType string, err error) {

	loginLog := model.LoginLog{}

	if vars.Config.LDAP != nil {
		err = vars.Config.LDAP.Authentication(cred.Username, cred.Password)
		if err == nil {
			authType = "ldap"
			loginLog.IPAddress = ctx.ClientIP()
			loginLog.TransactionAt = time.Now().In(vars.TZ)
			loginLog.SuccessLogin = true
			loginLog.ByUser = fmt.Sprintf("%s,%s", authType, cred.Username)
		}
	}

	if vars.Config.Radius != nil && loginLog.IPAddress == "" {
		err = vars.Config.Radius.Authentication(cred.Username, cred.Password)
		if err == nil {
			authType = "radius"
			loginLog.IPAddress = ctx.ClientIP()
			loginLog.TransactionAt = time.Now().In(vars.TZ)
			loginLog.SuccessLogin = true
			loginLog.ByUser = fmt.Sprintf("%s,%s", authType, cred.Username)
		}
	}

	if loginLog.IPAddress == "" {
		user := model.User{}

		if vars.Database.Where("username = ? and hashed = ?", cred.Username, utils.Sha512encode(cred.Password)).First(&user).Error != nil {
			loginLog.IPAddress = ctx.ClientIP()
			loginLog.TransactionAt = time.Now().In(vars.TZ)
			loginLog.SuccessLogin = false
			loginLog.ByUser = cred.Username

			if user.ID != "" {
				user.FailedLoginCount += 1
				vars.Database.Save(&user)
			}

			err = errors.New("your credential is not correct")

		} else {
			authType = "native"
			loginLog.IPAddress = ctx.ClientIP()
			loginLog.TransactionAt = time.Now().In(vars.TZ)
			loginLog.SuccessLogin = true
			loginLog.ByUser = fmt.Sprintf("%s,%s", authType, cred.Username)

			if user.ID != "" {
				user.FailedLoginCount = 0
				vars.Database.Save(&user)
			}

			err = nil
		}

		fmt.Println(user)
	}

	vars.Database.Create(&loginLog)

	return
}
