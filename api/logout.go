package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) logout(c *gin.Context) {

	clientIp := c.ClientIP()

	mSession := inmemory_model.Session{}
	err := vars.InMemoryDatabase.Where("ip_address = ?", clientIp).First(&mSession).Error
	if err != nil {
		msg := fmt.Sprintf("not found session of %s", clientIp)
		c.Redirect(http.StatusFound, fmt.Sprintf("/error?msg=%s", msg))
		return
	}

	err = session.CutOffSession(mSession)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), clientIp)
		c.Redirect(http.StatusFound, fmt.Sprintf("/error?msg=%s", msg))
		return
	}

	go func() {
		network.BWDel(&mSession)

		vars.Database.Create(&model.LogoutLog{
			TransactionAt: time.Now().In(vars.TZ),
			IPAddress:     clientIp,
			ByUser:        fmt.Sprintf("%s,%s", mSession.AuthType, mSession.User),
		})

	}()

	c.Redirect(http.StatusFound, "/login")
}

func (ctl *controller) logoutAllDevices(c *gin.Context) {

	clientIp := c.ClientIP()

	mSession := inmemory_model.Session{}
	err := vars.InMemoryDatabase.Where("ip_address = ?", clientIp).First(&mSession).Error
	if err != nil {
		msg := fmt.Sprintf("not found session of %s", clientIp)
		c.Redirect(http.StatusFound, fmt.Sprintf("/error?msg=%s", msg))
		return
	}

	err = session.CutOffSession(mSession)
	if err != nil {
		msg := fmt.Sprintf("%s via %s", err.Error(), clientIp)
		c.Redirect(http.StatusFound, fmt.Sprintf("/error?msg=%s", msg))
		return
	}

	go func() {
		network.BWDel(&mSession)

		vars.Database.Create(&model.LogoutLog{
			TransactionAt: time.Now().In(vars.TZ),
			IPAddress:     clientIp,
			ByUser:        fmt.Sprintf("%s,%s", mSession.AuthType, mSession.User),
		})

	}()

	allSession := []inmemory_model.Session{}
	vars.InMemoryDatabase.Where("user = ?", mSession.User).First(&allSession)

	for _, s := range allSession {
		session.CutOffSession(s)
		go func(s *inmemory_model.Session) {
			network.BWDel(s)

			vars.Database.Create(&model.LogoutLog{
				TransactionAt: time.Now().In(vars.TZ),
				IPAddress:     clientIp,
				ByUser:        fmt.Sprintf("%s,%s", mSession.AuthType, mSession.User),
			})

		}(&s)
	}

	c.Redirect(http.StatusFound, "/login")
}
