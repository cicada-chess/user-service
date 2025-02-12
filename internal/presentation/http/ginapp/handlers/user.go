package handlers

import (
	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
)

type UserHandler struct {
	Service interfaces.UserService
}

func (h *UserHandler) Ping(c *gin.Context) {
	c.JSON(200, "pong!")
}

func (h *UserHandler) Create(c *gin.Context) {
	user := &entity.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(400, err)
		return
	}
	id, err := h.Service.Create(user)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, id)
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.GetById(id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(200, user)
}
