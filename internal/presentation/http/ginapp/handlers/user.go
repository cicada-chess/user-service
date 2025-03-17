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

// @title UserSVC API
// @version 1.0
// @description API для управления пользователями

// @host localhost:8080
// @BasePath /api/v1

type UserHandler struct {
	Service interfaces.UserService
	Log     logrus.FieldLogger
}

func NewUserHandler(service interfaces.UserService, log logrus.FieldLogger) *UserHandler {
	return &UserHandler{
		Service: service,
		Log:     log,
	}
}

// CreateUser godoc
// @Summary Создание пользователя
// @Description Создаёт нового пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} response.SuccessResponse "Пользователь создан успешно"
// @Failure 400 {object} response.ErrorResponse "Ошибочные данные"
// @Failure 409 {object} response.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/create [post]
func (h *UserHandler) Create(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.Log.Errorf("Failed to bind user: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user := &entity.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
	createdUser, err := h.Service.Create(c.Request.Context(), user)
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
	response.NewSuccessResponse(c, http.StatusCreated, "Пользователь создан успешно", createdUser)

}

// GetUserById godoc
// @Summary Получение пользователя по ID
// @Description Возвращает данные пользователя по его идентификатору
// @Tags Users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} response.SuccessResponse "Данные пользователя найдены"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/{id} [get]
func (h *UserHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.GetById(c.Request.Context(), id)

	if err != nil {
		h.Log.Errorf("Failed to get user by id: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case application.ErrInvalidUUIDFormat:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	response.NewSuccessResponse(c, http.StatusOK, "Данные пользователя найдены успешно", user)

}

// UpdateUserInfo godoc
// @Summary Обновление данных пользователя
// @Description Изменяет информацию о пользователе (email, username и т.д.)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Param request body dto.UpdateInfoRequest true "Новые данные пользователя"
// @Success 200 {object} response.SuccessResponse "Обновление прошло успешно"
// @Failure 400 {object} response.ErrorResponse "Ошибочные данные"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/{id} [patch]
func (h *UserHandler) UpdateInfo(c *gin.Context) {
	id := c.Param("id")

	var request dto.UpdateInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.Log.Errorf("Failed to bind update info: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := h.Service.GetById(c.Request.Context(), id)
	if err != nil {
		h.Log.Errorf("Failed to get user by id: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case application.ErrInvalidUUIDFormat:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if request.Email != nil {
		userData.Email = *request.Email
	}
	if request.Username != nil {
		userData.Username = *request.Username
	}
	if request.Password != nil {
		userData.Password = *request.Password
	}
	if request.Rating != nil {
		userData.Rating = *request.Rating
	}
	if request.IsActive != nil {
		userData.IsActive = *request.IsActive
	}
	if request.Role != nil {
		userData.Role = *request.Role
	}

	updatedUser, err := h.Service.UpdateInfo(c.Request.Context(), userData)
	if err != nil {
		h.Log.Errorf("Failed to update user info: %v", err)
		switch err {
		case application.ErrInvalidIntegerValue:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверное числовое значение")
			return
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	response.NewSuccessResponse(c, http.StatusOK, "Информация о пользователе обновлена", updatedUser)
}

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаляет пользователя по ID
// @Tags Users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 204 {object} response.SuccessResponse "Пользователь удалён"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.Service.Delete(c.Request.Context(), id)
	if err != nil {
		h.Log.Errorf("Failed to delete user: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case application.ErrInvalidUUIDFormat:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusNoContent, "Пользователь удалён", nil)
}

// GetAllUsers godoc
// @Summary Получение списка пользователей
// @Description Возвращает список пользователей с поддержкой пагинации и сортировки
// @Tags Users
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество пользователей на странице"
// @Param search query string false "Строка поиска"
// @Param sort_by query string false "Поле для сортировки"
// @Param order query string false "Порядок сортировки (asc/desc)"
// @Success 200 {object} response.SuccessResponse "Список пользователей"
// @Failure 400 {object} response.ErrorResponse "Ошибочные параметры запроса"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users [get]
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

// ChangePassword godoc
// @Summary Изменение пароля
// @Description Изменяет пароль пользователя при наличии корректных старого пароля
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Param request body dto.ChangePasswordRequest true "Старый и новый пароль"
// @Success 200 {object} response.SuccessResponse "Пароль успешно изменён"
// @Failure 400 {object} response.ErrorResponse "Ошибочные данные"
// @Failure 401 {object} response.ErrorResponse "Неверный пароль"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/{id}/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id := c.Param("id")
	var request dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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
		case application.ErrInvalidUUIDFormat:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Пароль изменен успешно", nil)
}

// ToggleActive godoc
// @Summary Смена статуса активности
// @Description Переключает признак активности пользователя
// @Tags Users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} response.SuccessResponse "Статус изменён успешно"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/{id}/toggle-active [post]
func (h *UserHandler) ToggleActive(c *gin.Context) {
	id := c.Param("id")
	isActive, err := h.Service.ToggleActive(c.Request.Context(), id)
	if err != nil {
		h.Log.Errorf("Failed to toggle active: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case application.ErrInvalidUUIDFormat:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Статус активности пользователя изменен успешно", isActive)
}

// GetRating godoc
// @Summary Получение рейтинга
// @Description Возвращает текущий рейтинг пользователя
// @Tags Users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} response.SuccessResponse "Рейтинг получен"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/{id}/rating [get]
func (h *UserHandler) GetRating(c *gin.Context) {
	id := c.Param("id")
	rating, err := h.Service.GetRating(c.Request.Context(), id)
	if err != nil {
		h.Log.Errorf("Failed to get rating: %v", err)
		switch err {
		case application.ErrUserNotFound:
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case application.ErrInvalidUUIDFormat:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Рейтинг пользователя получен успешно", rating)
}

// UpdateRating godoc
// @Summary Обновление рейтинга
// @Description Увеличивает или уменьшает рейтинг пользователя на указанную величину
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Param request body dto.UpdateRatingRequest true "Изменение рейтинга"
// @Success 200 {object} response.SuccessResponse "Рейтинг успешно обновлён"
// @Failure 400 {object} response.ErrorResponse "Ошибочные данные"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка"
// @Router /users/{id}/update-rating [post]
func (h *UserHandler) UpdateRating(c *gin.Context) {
	id := c.Param("id")
	var request dto.UpdateRatingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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
		case application.ErrInvalidIntegerValue:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверное числовое значение")
			return
		case application.ErrInvalidUUIDFormat:
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.NewSuccessResponse(c, http.StatusOK, "Рейтинг пользователя обновлен успешно", rating)
}
