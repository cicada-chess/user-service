package ginapp

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.mai.ru/cicada-chess/backend/user-service/docs"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp/handlers"
)

func InitRoutes(r *gin.Engine, service interfaces.UserService, logger logrus.FieldLogger) {
	handler := handlers.NewUserHandler(service, logger)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/users")
	{
		api.POST("/create", handler.Create)
		api.GET("/:id", handler.GetById)
		api.PATCH("/:id", handler.UpdateInfo)
		api.DELETE("/:id", handler.Delete)
		api.GET("/", handler.GetAll)
		api.POST("/:id/change-password", handler.ChangePassword)
		api.POST("/:id/toggle-active", handler.ToggleActive)
		api.GET("/:id/rating", handler.GetRating)
		api.GET("/:id/update-rating", handler.UpdateRating)
	}
}
