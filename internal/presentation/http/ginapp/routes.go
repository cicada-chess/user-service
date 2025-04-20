package ginapp

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.mai.ru/cicada-chess/backend/user-service/docs"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp/handlers"
)

func InitRoutes(r *gin.Engine, service interfaces.UserService, logger logrus.FieldLogger) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"}, // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	handler := handlers.NewUserHandler(service, logger)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

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
		api.POST("/:id/update-rating", handler.UpdateRating)
	}

}
