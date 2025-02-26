package handlers

import (
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
	request := &entity.User{}
	if err := c.ShouldBindJSON(request); err != nil {
		h.Log.Errorf("Failed to bind user: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.Service.Create(c.Request.Context(), request)
	if err != nil {
		h.Log.Errorf("Failed to create user: %v", err)
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response.NewSuccessResponse(c, http.StatusCreated, "User successfully created", user)

}
func (h *UserHandler) GetInfo(c *gin.Context) {

}
func (h *UserHandler) UpdateInfo(c *gin.Context) {
}

func (h *UserHandler) Delete(c *gin.Context) {
}

func (h *UserHandler) GetAll(c *gin.Context) {
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
}

func (h *UserHandler) ToggleActive(c *gin.Context) {
}

func (h *UserHandler) GetRating(c *gin.Context) {
}

func (h *UserHandler) UpdateRating(c *gin.Context) {
}
