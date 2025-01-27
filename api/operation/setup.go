package api_operation

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/utils"
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
	envApiKeyHashed := os.Getenv("COCONG_API_KEY_HASHED")

	if utils.Sha512encode(tokenString) == envApiKeyHashed {
		c.Next()
		return
	}

	_, ok := vars.AdminSession.Get(tokenString)
	if !ok {
		c.JSON(http.StatusUnauthorized, "not authorised")
		c.Abort()
		return
	}

	c.Next()
}

func tokenMiddlewareSessionOnly(c *gin.Context) {
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
	ctl.router.GET("/op/me", tokenMiddlewareSessionOnly, ctl.me)
	ctl.router.POST("/op/change-password", tokenMiddlewareSessionOnly, ctl.adminChangePassword)

	session := ctl.router.Group("/op/session")
	session.GET("/query", tokenMiddleware, ctl.sessionQuery)
	session.PATCH("/kick", tokenMiddlewareSessionOnly, ctl.sessionKick)

	bw := ctl.router.Group("/op/bandwidth")
	bw.Use(tokenMiddlewareSessionOnly)

	bw.GET("/query", ctl.bandwidthQuery)
	bw.POST("/create", ctl.bandwidthCreate)
	bw.DELETE("/delete", ctl.bandwidthDelete)

	directory := ctl.router.Group("/op/directory")
	directory.Use(tokenMiddlewareSessionOnly)

	directory.GET("/query", ctl.directoryQuery)
	directory.POST("/create", ctl.directoryCreate)
	directory.PUT("/update", ctl.directoryUpdate)
	directory.DELETE("/delete", ctl.directoryDelete)

	user := ctl.router.Group("/op/user")
	user.Use(tokenMiddlewareSessionOnly)

	user.GET("/query", ctl.userQuery)
	user.POST("/create", ctl.userCreate)
	user.PUT("/update", ctl.userUpdate)
	user.DELETE("/delete", ctl.userDelete)
	user.PATCH("/password-reset", ctl.userPasswordReset)

	administrator := ctl.router.Group("/op/administrator")
	administrator.Use(tokenMiddlewareSessionOnly)

	administrator.GET("/query", ctl.administratorQuery)
	administrator.POST("/create", ctl.administratorCreate)
	administrator.PUT("/update", ctl.administratorUpdate)
	administrator.DELETE("/delete", ctl.administratorDelete)
	administrator.PATCH("/password-reset", ctl.administratorPasswordReset)

	logs := ctl.router.Group("/op/log")
	logs.Use(tokenMiddlewareSessionOnly)

	logs.GET("/login", ctl.loginLogQuery)
	logs.GET("/logout", ctl.logoutLogQuery)
	logs.GET("/net", ctl.networkLogQuery)
	logs.GET("/login-dump", ctl.loginLogDump)
	logs.GET("/logout-dump", ctl.logoutLogDump)

	sys := ctl.router.Group("/op/system")
	sys.Use(tokenMiddlewareSessionOnly)

	sys.GET("/config", ctl.getConfig)
	sys.PUT("/config", ctl.writeConfig)

	sys.PATCH("/service/core/restart", ctl.restartService)
	sys.PATCH("/service/dns/restart", ctl.restartDNSServer)
}
