package ginapp

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp/handlers"
)

func InitRoutes(r *gin.Engine, service interfaces.UserService, logger logrus.FieldLogger) {
	handler := &handlers.UserHandler{Service: service, Log: logger}
	r.POST("/ping", handler.Ping)

	api := r.Group("/users")
	{
		api.POST("/create", handler.Create)
		api.GET("/:id", handler.GetInfo)
	}
}
