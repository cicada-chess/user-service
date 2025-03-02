package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/user"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/user/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/response"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp/dto"
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
		switch err {
		case application.ErrEmailExists:
			response.NewErrorResponse(c, http.StatusConflict, "Данный email уже зарегистрирован")
			return
		case application.ErrUsernameExists:
			response.NewErrorResponse(c, http.StatusConflict, "Данный username уже зарегистрирован")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	response.NewSuccessResponse(c, http.StatusCreated, "Пользователь создан успешно", user)

}

func (h *UserHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.GetById(c.Request.Context(), id)

	if err != nil {
		h.Log.Errorf("Failed to get user by id: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	response.NewSuccessResponse(c, http.StatusOK, "Данные пользователя найдены успешно", user)

}

func (h *UserHandler) UpdateInfo(c *gin.Context) {
	id := c.Param("id")
	request := make(map[string]interface{})
	if err := c.ShouldBindJSON(request); err != nil {
		h.Log.Errorf("Failed to bind user: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user := &entity.User{ID: id}

	if email, exist := request["email"]; exist {
		user.Email = email.(string)
	}

	if username, exist := request["username"]; exist {
		user.Username = username.(string)
	}

	if password, exist := request["password"]; exist {
		user.Password = password.(string)
	}

	if role, exist := request["role"]; exist {
		user.Role = int(role.(float64))
	}

	if rating, exist := request["rating"]; exist {
		user.Rating = int(rating.(float64))
	}

	if isActive, exist := request["is_active"]; exist {
		user.IsActive = isActive.(bool)
	}

	updatedUser, err := h.Service.UpdateInfo(c.Request.Context(), user)
	if err != nil {
		h.Log.Errorf("Failed to update user: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case entity.ErrPasswordTooShort:
			response.NewErrorResponse(c, http.StatusBadRequest, "Пароль слишком короткий")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	response.NewSuccessResponse(c, http.StatusOK, "Данные пользователя обновлены успешно", updatedUser)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.Service.Delete(c.Request.Context(), id)
	if err != nil {
		h.Log.Errorf("Failed to delete user: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusNoContent, "Пользователь удален успешно", nil)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	var request dto.GetAllUsersRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		h.Log.Errorf("Failed to bind query parameters: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.Service.GetAll(c.Request.Context(), request.Page, request.Limit, request.Search, request.SortBy, request.Order)
	if err != nil {
		h.Log.Errorf("Failed to get all users: %v", err)
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.NewSuccessResponse(c, http.StatusOK, "Пользователи получены успешно", users)

}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")
	request := &dto.ChangePasswordRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		h.Log.Errorf("Failed to bind request: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Service.ChangePassword(c.Request.Context(), id, request.OldPassword, request.NewPassword)
	if err != nil {
		h.Log.Errorf("Failed to change password: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case application.ErrInvalidPassword:
			response.NewErrorResponse(c, http.StatusUnauthorized, "Неверный пароль")
			return
		case entity.ErrPasswordTooShort:
			response.NewErrorResponse(c, http.StatusBadRequest, "Новый пароль слишком короткий")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Пароль изменен успешно", nil)
}

func (h *UserHandler) ToggleActive(c *gin.Context) {
	id := c.Param("id")
	isActive, err := h.Service.ToggleActive(c.Request.Context(), id)
	if err != nil {
		h.Log.Errorf("Failed to toggle active: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Статус активности пользователя изменен успешно", isActive)
}

func (h *UserHandler) GetRating(c *gin.Context) {
	id := c.Param("id")
	rating, err := h.Service.GetRating(c.Request.Context(), id)
	if err != nil {
		h.Log.Errorf("Failed to get rating: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Рейтинг пользователя получен успешно", rating)
}

func (h *UserHandler) UpdateRating(c *gin.Context) {
	id := c.Param("id")
	request := &dto.UpdateRatingRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		h.Log.Errorf("Failed to bind request: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rating, err := h.Service.UpdateRating(c.Request.Context(), id, request.Delta)
	if err != nil {
		h.Log.Errorf("Failed to update rating: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Рейтинг пользователя обновлен успешно", rating)
}
