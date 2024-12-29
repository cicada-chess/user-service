package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/web"
)

func main() {
	app := gin.New()
	web.InitRoutes(app)

	app.Run(":8000")
}
