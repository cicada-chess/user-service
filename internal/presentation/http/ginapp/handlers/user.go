package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/response"
)

type UserHandler struct {
	Service interfaces.UserService
}

func (h *UserHandler) Ping(c *gin.Context) {
	response.NewSuccessResponse(c, http.StatusOK, "pong", nil)
}

func (h *UserHandler) Create(c *gin.Context) {
	user := &entity.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.Service.Create(user)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.NewSuccessResponse(c, http.StatusOK, fmt.Sprintf("User created with id: %s", id), nil)
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.GetById(id)
	if user == nil {
		response.NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.NewSuccessResponse(c, http.StatusOK, "success", user)
}
