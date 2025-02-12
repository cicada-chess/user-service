package ginapp

import (
	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp/handlers"
)

func InitRoutes(r *gin.Engine, service interfaces.UserService) {
	handler := &handlers.UserHandler{Service: service}
	r.POST("/ping", handler.Ping)

	api := r.Group("/users")
	{
		api.POST("/create", handler.Create)
		api.GET("/:id", handler.GetInfo)
	}
}
