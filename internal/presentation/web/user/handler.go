package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/web/response"
)

type Handler struct {
}

func (h *Handler) Ping(c *gin.Context) {
	response.NewSuccessResponse(c, http.StatusOK, "pong!")
}
