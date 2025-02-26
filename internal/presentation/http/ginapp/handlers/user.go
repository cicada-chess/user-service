package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
)

type UserHandler struct {
	Service interfaces.UserService
	Log     logrus.FieldLogger
}

func (h *UserHandler) Create(c *gin.Context) {

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
