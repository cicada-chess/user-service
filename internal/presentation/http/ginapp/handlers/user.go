package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/response"
)

type UserHandler struct {
	Service interfaces.UserService
	Log     logrus.FieldLogger
}

func (h *UserHandler) Create(c *gin.Context) {
	user := &entity.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		h.Log.Errorf("Failed to bind user: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.Service.Create(c, user)
	if err != nil {
		h.Log.Errorf("Failed to create user: %v", err)
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.NewSuccessResponse(c, http.StatusOK, fmt.Sprintf("User created with id: %s", id), nil)
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.GetById(c, id)
	if user == nil {
		h.Log.Errorf("User not found: %s", id)
		response.NewErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		h.Log.Errorf("Failed to get user: %v", err)
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.NewSuccessResponse(c, http.StatusOK, "success", user)
}
