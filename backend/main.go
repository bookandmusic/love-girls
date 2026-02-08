package main

import (
	"errors"
	"log"
	"os"

	"github.com/bookandmusic/love-girl/docs"
	"github.com/bookandmusic/love-girl/internal/config"
	"github.com/bookandmusic/love-girl/internal/server"
	"github.com/bookandmusic/love-girl/provider"
)

var (
	Version   = "0.1.0"
	Commit    = "unknown"
	BuildTime = "unknown"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.oauth2.password     OAuth2Password
//	@tokenUrl								/api/v1/user/token
//	@scope.read								Grants read access
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

func main() {
	restartCh := make(chan struct{}, 1)

	for {
		app, cleanup, err := provider.InitApp()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		docs.SwaggerInfo.Host = app.Server.HostName

		err = app.Run(Version, Commit, BuildTime, restartCh)

		if errors.Is(err, server.ErrRestart) {
			log.Println("配置变更，正在重启服务...")
			cleanup()
			config.ResetConfig()
			continue
		}

		if err != nil {
			os.Exit(1)
		}
		break
	}
}
