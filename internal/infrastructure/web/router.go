package web

import (
	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/web/user"
)

func InitRoutes(engine *gin.Engine) {
	handler := user.Handler{}
	api := engine.Group("/api")
	{
		api.GET("/ping", handler.Ping)
	}
}
