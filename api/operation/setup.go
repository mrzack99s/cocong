package api_operation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/vars"
)

type controller struct {
	router gin.IRouter
}

func NewController(router gin.IRouter) *controller {
	s := &controller{
		router: router,
	}
	s.register()
	return s
}

func tokenMiddleware(c *gin.Context) {
	tokenString := c.Request.Header.Get("api-token")

	_, ok := vars.AdminSession.Get(tokenString)
	if !ok {
		c.JSON(http.StatusUnauthorized, "not authorised")
		c.Abort()
		return
	}
	c.Next()
}

func (ctl *controller) register() {

	ctl.router.POST("/op/login", ctl.login)
	ctl.router.DELETE("/op/logout", ctl.logout)
	ctl.router.POST("/op/refresh-token", ctl.refreshToken)
	ctl.router.GET("/op/me", tokenMiddleware, ctl.me)
	ctl.router.POST("/op/change-password", tokenMiddleware, ctl.adminChangePassword)

	session := ctl.router.Group("/op/session")
	session.Use(tokenMiddleware)

	session.GET("/query", ctl.sessionQuery)
	session.PATCH("/kick", ctl.sessionKick)

	bw := ctl.router.Group("/op/bandwidth")
	bw.Use(tokenMiddleware)

	bw.GET("/query", ctl.bandwidthQuery)
	bw.POST("/create", ctl.bandwidthCreate)
	bw.DELETE("/delete", ctl.bandwidthDelete)

	directory := ctl.router.Group("/op/directory")
	directory.Use(tokenMiddleware)

	directory.GET("/query", ctl.directoryQuery)
	directory.POST("/create", ctl.directoryCreate)
	directory.PUT("/update", ctl.directoryUpdate)
	directory.DELETE("/delete", ctl.directoryDelete)

	user := ctl.router.Group("/op/user")
	user.Use(tokenMiddleware)

	user.GET("/query", ctl.userQuery)
	user.POST("/create", ctl.userCreate)
	user.PUT("/update", ctl.userUpdate)
	user.DELETE("/delete", ctl.userDelete)
	user.PATCH("/password-reset", ctl.userPasswordReset)

	administrator := ctl.router.Group("/op/administrator")
	administrator.Use(tokenMiddleware)

	administrator.GET("/query", ctl.administratorQuery)
	administrator.POST("/create", ctl.administratorCreate)
	administrator.PUT("/update", ctl.administratorUpdate)
	administrator.DELETE("/delete", ctl.administratorDelete)
	administrator.PATCH("/password-reset", ctl.administratorPasswordReset)

	logs := ctl.router.Group("/op/log")
	logs.Use(tokenMiddleware)

	logs.GET("/login", ctl.loginLogQuery)
	logs.GET("/logout", ctl.logoutLogQuery)
	logs.GET("/net", ctl.networkLogQuery)
	logs.GET("/login-dump", ctl.loginLogDump)
	logs.GET("/logout-dump", ctl.logoutLogDump)
	logs.GET("/net-dump", ctl.networkLogDump)
}
