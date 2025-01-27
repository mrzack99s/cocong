package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/services"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) getAuthentication(c *gin.Context) {

	clientIp := c.ClientIP()

	username := c.PostForm("username")
	password := c.PostForm("password")

	checkCredential := types.CredentialVerification{
		Username: username,
		Password: password,
	}

	newSession := types.SessionInfo{
		IPAddress: clientIp,
		User:      checkCredential.Username,
	}

	authType, err := services.Authentication(c, checkCredential)
	if err != nil {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", err.Error()))
		return
	}

	newSession.AuthType = authType

	countLoggedUsername := 0
	listUsernameSession, err := session.Instance.GetByUsername(checkCredential.Username)
	if err == nil {
		countLoggedUsername = len(listUsernameSession)
	}

	if newSession.AuthType == "native" {
		var user model.User
		vars.Database.Preload("Directory").Preload("Directory.Bandwidth").Where("username = ?", checkCredential.Username).First(&user)

		if !user.Directory.Enable {
			msg := "your directory is disabled"
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", msg))
			return
		}

		if !user.Enable {
			msg := "your account is disabled"
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", msg))
			return
		}

		newSession.Bandwidth = *user.Directory.Bandwidth

		if user.Directory.MaxConcurrent > 0 && int64(countLoggedUsername) >= user.Directory.MaxConcurrent {
			msg := fmt.Sprintf("user %s reached the limit concurrent session and login via %s", checkCredential.Username, clientIp)
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", msg))
			return
		}

		if user.Directory.Bandwidth != nil {
			network.BWSet(&newSession)
		}

	} else {
		if vars.Config.MaxConcurrentSession > 0 && int64(countLoggedUsername) >= int64(vars.Config.MaxConcurrentSession) {
			msg := fmt.Sprintf("user %s reached the limit concurrent session and login via %s", checkCredential.Username, clientIp)
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", msg))
			return
		}
	}

	err = session.Instance.Create(newSession)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), clientIp)
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", msg))
		return
	}

	c.Redirect(http.StatusSeeOther, "/status")
}

func (ctl *controller) changePassword(c *gin.Context) {

	currentPassword := c.PostForm("current_password")
	newPassword := c.PostForm("new_password")
	newPasswordAgain := c.PostForm("new_again_password")

	clientIp := c.ClientIP()

	// session := inmemory_model.Session{}
	session, err := session.Instance.GetByIP(clientIp)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	var user model.User
	err = vars.Database.Where("username = ?", session.User).First(&user).Error
	if err != nil {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", "not found this username, change password only native authentication support"))
		return
	}

	if utils.Sha512encode(currentPassword) != user.Hashed {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", "current password is not correct"))
		return
	}

	if utils.Sha512encode(newPassword) == user.Hashed {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", "your new password is current password. please change to new password"))
		return
	}

	if newPasswordAgain != newPassword {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", "cannot change password, new password mismatch"))
		return
	}

	user.Hashed = utils.Sha512encode(newPassword)
	err = vars.Database.Save(&user).Error
	if err != nil {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/error?msg=%s", "cannot change password, please contact your administrator"))
		return
	}

	c.Redirect(http.StatusSeeOther, "/status")
}
