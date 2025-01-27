package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/cocong/api"
	api_operation "github.com/mrzack99s/cocong/api/operation"
	"github.com/mrzack99s/cocong/constants"
	"github.com/mrzack99s/cocong/network"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/setup"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
	"github.com/mrzack99s/cocong/watcher"
	"github.com/spf13/cobra"
)

func main() {

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "To run an application. Default run in development mode",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			if vars.SYS_DEBUG {
				vars.Config.TimeZone = "Asia/Bangkok"
				vars.Config.AuthorizedNetworks = []string{"0.0.0.0/0", "::/0"}
				network.AuthorizedNetworks = vars.Config.AuthorizedNetworks
				setup.AppConfig()

			} else {
				if !utils.IsRootPrivilege() {
					panic(`this application needs the ability to run commands as root. We are unable to find either "sudo" or "su" available to make this happen.`)
				}

				if _, err := os.Stat(constants.LOG_DIR); os.IsNotExist(err) {
					os.Mkdir(constants.LOG_DIR, 0644)
				}

				appLogger := setup.LoggingConfig{
					Directory:  constants.LOG_DIR,
					Filename:   "applog",
					MaxSize:    50,
					MaxAge:     90,
					MaxBackups: 90,
				}
				appLogger.Configure()
				setup.AppConfig()

			}

			utils.VerifyTimeZone()
			setup.GetDeviceResources()
			setup.Database()
			session.Instance.New()

			err := vars.Config.LDAP.NewLDAPConnectionPool()
			if err != nil {
				panic(err)
			}

			if !vars.SYS_DEBUG {
				gin.SetMode(gin.ReleaseMode)

				err := network.InitializeCaptivePortal()
				if err != nil {
					panic(err)
				}

				watcher.CaptivePortalDetector(context.Background())
				watcher.NetIdleChecking(context.Background())
				watcher.NetWatcher(context.Background())
				watcher.ConntrackChecking(context.Background())

			}

			router := gin.Default()

			corsConfig := cors.DefaultConfig()
			corsConfig.AllowAllOrigins = true
			corsConfig.AllowHeaders = []string{"Content-Type, Content-Length, Accept-Encoding, origin, Cache-Control, api-token, access-control-allow-origin"}

			corsMiddleware := cors.New(corsConfig)

			router.Use(corsMiddleware)

			if vars.SYS_DEBUG {
				router.LoadHTMLGlob("templates/*")
			} else {
				router.LoadHTMLGlob(constants.APP_DIR + "/templates/*")
			}

			api.Newcontroller(router)
			api_operation.NewController(router)
			router.NoRoute(func(ctx *gin.Context) {
				ctx.Data(404, "text/plain", nil)
				ctx.Abort()
			})

			if vars.SYS_DEBUG {
				// listener, err := net.Listen("tcp4", "0.0.0.0:4443")
				// if err != nil {
				// 	log.Fatalf("Failed to create listener: %v", err)
				// }

				// server := &http.Server{
				// 	Handler:      router,
				// 	ReadTimeout:  10 * time.Second,
				// 	WriteTimeout: 10 * time.Second,
				// }

				// if err := server.ServeTLS(listener, "./certs/server.crt", "./certs/server.key"); err != nil {
				// 	log.Fatalf("Server failed: %v", err)
				// }
				router.RunTLS("0.0.0.0:4443", "./certs/server.crt", "./certs/server.key")

			} else {
				router.RunTLS("0.0.0.0:443", constants.CONFIG_DIR+"/certs/server.crt", constants.CONFIG_DIR+"/certs/server.key")
			}

		},
	}
	cmdRun.Flags().BoolVarP(&vars.SYS_DEBUG, "debug", "d", false, "Run in debug mode")

	var cmdCertificate = &cobra.Command{
		Use:   "gencert",
		Short: "To generate a self-signed certificate.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Println("usage: cocong gencert [domain-name]")
				os.Exit(1)
			}

			vars.Config.DomainName = args[0]
			utils.GenerateSelfSignCert()
		},
	}

	var showVersion = &cobra.Command{
		Use:   "version",
		Short: "To show version",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CoCoNG version " + constants.VERSION)
		},
	}

	var rootCmd = &cobra.Command{Use: "cocong"}
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdCertificate)
	rootCmd.AddCommand(showVersion)
	rootCmd.Execute()
}
