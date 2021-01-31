package Router

import (
	"Mmx/Controllers"
	"Mmx/Modules"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(gin.ReleaseMode)
	G := gin.New()
	G.Use(gin.Recovery())

	G.GET(Modules.Config.Path, Controllers.Main)

	Modules.Callback.CMD.Info("开始监听 0.0.0.0:" + Modules.Config.Port + Modules.Config.Path)
	if err := G.Run(":" + Modules.Config.Port); err != nil {
		Modules.Callback.CMD.Error(err)
	}
}
