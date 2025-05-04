package ginapp

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gitlab.mai.ru/cicada-chess/backend/user-service/docs"
	profileInterfaces "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
	userInterfaces "gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp/handlers"
)

func InitRoutes(
	r *gin.Engine,
	userService userInterfaces.UserService,
	profileService profileInterfaces.ProfileService,
	logger logrus.FieldLogger,
) {
	r.Static("/uploads/avatars", "/uploads/avatars")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://cikada-inky.vercel.app"}, // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path != "/health" {
			logger.Infof("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		}
		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := handlers.NewUserHandler(userService, logger)
	profileHandler := handlers.NewProfileHandler(profileService, logger)

	users := r.Group("/users")
	{
		users.POST("/create", userHandler.Create)
		users.GET("/:id", userHandler.GetById)
		users.PATCH("/:id", userHandler.UpdateInfo)
		users.DELETE("/:id", userHandler.Delete)
		users.GET("", userHandler.GetAll)
		users.POST("/:id/change-password", userHandler.ChangePassword)
		users.POST("/:id/toggle-active", userHandler.ToggleActive)
		users.GET("/:id/rating", userHandler.GetRating)
		users.POST("/:id/update-rating", userHandler.UpdateRating)
	}

	profile := r.Group("/profile")
	{
		profile.POST("/create/:id", profileHandler.CreateProfile)
		profile.GET("", profileHandler.GetProfile)
		profile.PATCH("", profileHandler.UpdateProfile)
		profile.POST("/avatar", profileHandler.UploadAvatar)
	}

}
