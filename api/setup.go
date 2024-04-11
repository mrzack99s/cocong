package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

type controller struct {
	router gin.IRouter
}

func Newcontroller(router gin.IRouter) *controller {
	s := &controller{
		router: router,
	}
	s.register()
	return s
}

func GetUnAuthirizedNetworkMiddleware(c *gin.Context) {

	path := c.Request.URL.Path
	if path == "/unauthorised" {
		c.Next()
		return
	}

	clientIp := c.ClientIP()

	allow := false
	for _, cidr := range vars.Config.AuthorizedNetworks {
		if utils.Ipv4InCidr(cidr, clientIp) {
			allow = true
			break
		}
	}

	if !allow {
		c.Redirect(http.StatusFound, "/unauthorised")
		c.Abort()
		return
	}

	c.Next()

}

func (ctl *controller) register() {

	ctl.router.Any(LoginEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.authLoginPage)
	ctl.router.Any(LogoutEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.logout)
	ctl.router.Any(LogoutAllDeviceEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.logoutAllDevices)

	ctl.router.Any(StatusEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.getStatus)
	ctl.router.Any(ErrorEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.authErrorPage)
	ctl.router.Any(ChangePasswordPageEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.changePasswordPage)

	ctl.router.POST(ChangePasswordEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.changePassword)
	ctl.router.POST(AuthEndpoint, GetUnAuthirizedNetworkMiddleware, ctl.getAuthentication)

	ctl.router.Any(UnauthorisedEndpoint, ctl.unauthorised)
}
