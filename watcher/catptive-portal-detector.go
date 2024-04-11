package watcher

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/constants"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
)

func CaptivePortalDetector(ctx context.Context) {
	go func(ctx context.Context) {

		intIp, err := utils.GetSecureInterfaceIpv4Addr()
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				if !vars.SYS_DEBUG {
					gin.SetMode(gin.ReleaseMode)
				}

				router := gin.Default()
				router.Use(cors.Default())
				router.Any("/", func(c *gin.Context) {
					if utils.StringContains(vars.URL_CAPTIVE_PORTAL_DETECTION, c.Request.Host+c.Request.URL.Path) {
						if vars.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, vars.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)
						}
					} else {
						if vars.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, vars.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)
						}
					}
				})
				router.NoRoute(func(c *gin.Context) {
					if vars.Config.ExternalPortalURL != "" {
						c.Redirect(http.StatusFound, vars.Config.ExternalPortalURL)
						return
					} else {
						redirect(c, intIp)
					}
				})
				err := router.Run(fmt.Sprintf("%s:8080", intIp))
				if err != nil {
					vars.SystemLog.Println("captive-portal-detect-http: " + err.Error())
					return
				}
			}
		}
	}(ctx)

	go func(ctx context.Context) {

		intIp, err := utils.GetSecureInterfaceIpv4Addr()
		if err != nil {
			panic(err)
		}

		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				if !vars.SYS_DEBUG {
					gin.SetMode(gin.ReleaseMode)
				}

				router := gin.Default()
				router.Use(cors.Default())
				router.Any("/", func(c *gin.Context) {
					if utils.StringContains(vars.URL_CAPTIVE_PORTAL_DETECTION, c.Request.Host+c.Request.URL.Path) {
						if vars.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, vars.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)

						}
					} else {
						if vars.Config.ExternalPortalURL != "" {
							c.Redirect(http.StatusFound, vars.Config.ExternalPortalURL)
							return
						} else {
							redirect(c, intIp)
						}
					}
				})
				router.NoRoute(func(c *gin.Context) {
					if vars.Config.ExternalPortalURL != "" {
						c.Redirect(http.StatusFound, vars.Config.ExternalPortalURL)
						return
					} else {
						redirect(c, intIp)
					}
				})

				err := router.RunTLS(fmt.Sprintf("%s:8443", intIp), constants.CONFIG_DIR+"/certs/server.crt", constants.CONFIG_DIR+"/certs/server.key")
				if err != nil {
					err := router.RunTLS(fmt.Sprintf("%s:8443", intIp), "./certs/server.crt", "./certs/server.key")
					if err != nil {
						vars.SystemLog.Println("captive-portal-detect-https: " + err.Error())
						return
					}
				}
			}
		}
	}(ctx)
}

func redirect(c *gin.Context, intIp string) {

	if vars.Config.DomainName != "" {
		c.Redirect(http.StatusFound, fmt.Sprintf("https://%s/login", vars.Config.DomainName))
		return
	} else {
		c.Redirect(http.StatusFound, fmt.Sprintf("https://%s/login", intIp))
		return
	}

}
