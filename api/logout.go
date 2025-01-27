package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/types"
	"github.com/mrzack99s/cocong/vars"
)

func (ctl *controller) logout(c *gin.Context) {

	clientIp := c.ClientIP()

	mSession, err := session.Instance.GetByIP(clientIp)

	if err != nil {
		msg := fmt.Sprintf("not found session of %s", clientIp)
		c.Redirect(http.StatusFound, fmt.Sprintf("/error?msg=%s", msg))
		return
	}

	session.Instance.Delete(mSession.IPAddress)

	go func() {
		network.BWDel(mSession)

		vars.Database.Create(&model.LogoutLog{
			TransactionAt: time.Now().In(vars.TZ),
			IPAddress:     clientIp,
			User:          fmt.Sprintf("%s,%s", mSession.AuthType, mSession.User),
		})

	}()

	c.Redirect(http.StatusFound, "/login")
}

func (ctl *controller) logoutAllDevices(c *gin.Context) {

	clientIp := c.ClientIP()

	// mSession, err := utils.RedisGetInsideWildcard[inmemory_model.Session](context.Background(), vars.RedisCache, fmt.Sprintf("session|*|%s", clientIp))
	// if err != nil {
	// 	msg := fmt.Sprintf("not found session of %s", clientIp)
	// 	c.Redirect(http.StatusFound, fmt.Sprintf("/error?msg=%s", msg))
	// 	return
	// }

	// session.CutOffSession(mSession)

	// go func() {
	// 	network.BWDel(&mSession)

	// 	vars.Database.Create(&model.LogoutLog{
	// 		TransactionAt: time.Now().In(vars.TZ),
	// 		IPAddress:     clientIp,
	// 		User:          fmt.Sprintf("%s,%s", mSession.AuthType, mSession.User),
	// 	})

	// }()

	mSession, err := session.Instance.GetByIP(clientIp)
	if err != nil {
		msg := fmt.Sprintf("not found session of %s", clientIp)
		c.Redirect(http.StatusFound, fmt.Sprintf("/error?msg=%s", msg))
		return
	}

	mSessions, _ := session.Instance.GetByUsername(mSession.User)
	for _, s := range mSessions {
		session.Instance.Delete(s.IPAddress)
		go func(s *types.SessionInfo) {
			network.BWDel(s)

			vars.Database.Create(&model.LogoutLog{
				TransactionAt: time.Now().In(vars.TZ),
				IPAddress:     clientIp,
				User:          fmt.Sprintf("%s,%s", mSession.AuthType, mSession.User),
			})

		}(&s)
	}

	c.Redirect(http.StatusFound, "/login")
}
