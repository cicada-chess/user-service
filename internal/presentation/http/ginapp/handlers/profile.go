package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "gitlab.mai.ru/cicada-chess/backend/user-service/docs"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/interfaces"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/infrastructure/response"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/presentation/http/ginapp/dto"
)

type ProfileHandler struct {
	profileService interfaces.ProfileService
	logger         logrus.FieldLogger
}

func NewProfileHandler(profileService interfaces.ProfileService, logger logrus.FieldLogger) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger,
	}
}

// CreateProfile godoc
// @Summary Создание профиля пользователя
// @Description Создает профиль для пользователя по его идентификатору
// @Tags Profile
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} docs.SuccessResponse{data=docs.Profile} "Профиль создан"
// @Failure 404 {object} docs.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} docs.ErrorResponse "Внутренняя ошибка сервера"
// @Router /profile/create/{id} [post]
func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	id := c.Param("id")

	profile, err := h.profileService.CreateProfile(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorf("failed to create profile for user: %v", err)
		switch {
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	profileDTO := &dto.Profile{
		UserID:      profile.UserID,
		Username:    profile.Username,
		Age:         profile.Age,
		Location:    profile.Location,
		Description: profile.Description,
		AvatarURL:   profile.AvatarURL,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}

	response.NewSuccessResponse(c, http.StatusOK, "Профиль создан успешно", profileDTO)
}

// GetProfile godoc
// @Summary Получение профиля пользователя
// @Description Возвращает профиль текущего аутентифицированного пользователя
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} docs.SuccessResponse{data=docs.Profile} "Профиль пользователя"
// @Failure 401 {object} docs.ErrorResponse "Ошибка авторизации"
// @Failure 500 {object} docs.ErrorResponse "Внутренняя ошибка сервера"
// @Router /profile [get]
// @Security BearerAuth
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		response.NewErrorResponse(c, http.StatusUnauthorized, "Неавторизованный доступ")
		return
	}

	userID, err := h.profileService.GetUserIDFromToken(c, tokenHeader)
	if err != nil {
		h.logger.Errorf("Failed to get user id from token: %v", err)
		switch {
		case errors.Is(err, application.ErrTokenInvalidOrExpired):
			response.NewErrorResponse(c, http.StatusUnauthorized, "Токен недействителен или истёк")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get profile: %v", err)
		switch {
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case errors.Is(err, application.ErrProfileNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Профиль пользователя не найден")
			return
		case errors.Is(err, application.ErrInvalidUUIDFormat):
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	profileDTO := &dto.Profile{
		UserID:      profile.UserID,
		Username:    profile.Username,
		Description: profile.Description,
		Age:         profile.Age,
		Location:    profile.Location,
		AvatarURL:   profile.AvatarURL,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}

	response.NewSuccessResponse(c, http.StatusOK, "Профиль получен успешно", profileDTO)
}

// UpdateProfile godoc
// @Summary Обновление профиля пользователя
// @Description Обновляет профиль текущего аутентифицированного пользователя
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body docs.UpdateProfileRequest true "Данные профиля"
// @Success 200 {object} docs.SuccessResponse{data=docs.Profile} "Профиль обновлен"
// @Failure 400 {object} docs.ErrorResponse "Неверные данные профиля"
// @Failure 401 {object} docs.ErrorResponse "Ошибка авторизации"
// @Failure 500 {object} docs.ErrorResponse "Внутренняя ошибка сервера"
// @Router /profile [patch]
// @Security BearerAuth
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		response.NewErrorResponse(c, http.StatusUnauthorized, "Неавторизованный доступ")
		return
	}

	userID, err := h.profileService.GetUserIDFromToken(c, tokenHeader)
	if err != nil {
		h.logger.Errorf("failed to get user id from token: %v", err)
		switch {
		case errors.Is(err, application.ErrTokenInvalidOrExpired):
			response.NewErrorResponse(c, http.StatusUnauthorized, "Токен недействителен или истёк")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	var request dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Errorf("Failed to bind profile update request: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get profile: %v", err)
		switch {
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case errors.Is(err, application.ErrProfileNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Профиль пользователя не найден")
			return
		case errors.Is(err, application.ErrInvalidUUIDFormat):
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if request.Description != nil {
		profile.Description = *request.Description
	}
	if request.Age != nil {
		profile.Age = *request.Age
	}
	if request.Location != nil {
		profile.Location = *request.Location
	}
	if request.AvatarURL != nil {
		profile.AvatarURL = *request.AvatarURL
	}

	updatedProfile, err := h.profileService.UpdateProfile(c.Request.Context(), profile)
	if err != nil {
		h.logger.Errorf("Failed to update profile: %v", err)
		switch {
		case errors.Is(err, application.ErrInvalidAge):
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный возраст")
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	profileDTO := &dto.Profile{
		UserID:      updatedProfile.UserID,
		Username:    updatedProfile.Username,
		Description: updatedProfile.Description,
		Age:         updatedProfile.Age,
		Location:    updatedProfile.Location,
		AvatarURL:   updatedProfile.AvatarURL,
		CreatedAt:   updatedProfile.CreatedAt,
		UpdatedAt:   updatedProfile.UpdatedAt,
	}

	response.NewSuccessResponse(c, http.StatusOK, "Профиль обновлен успешно", profileDTO)
}

// UploadAvatar godoc
// @Summary Загрузка аватара
// @Description Загружает аватар для текущего аутентифицированного пользователя
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param avatar formData file true "Файл аватара"
// @Success 200 {object} docs.SuccessResponse{data=docs.Profile} "Аватар загружен"
// @Failure 400 {object} docs.ErrorResponse "Ошибка загрузки файла"
// @Failure 401 {object} docs.ErrorResponse "Ошибка авторизации"
// @Failure 500 {object} docs.ErrorResponse "Внутренняя ошибка сервера"
// @Router /profile/avatar [post]
// @Security BearerAuth
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	tokenHeader := c.GetHeader("Authorization")
	if tokenHeader == "" {
		response.NewErrorResponse(c, http.StatusUnauthorized, "Неавторизованный доступ")
		return
	}

	userID, err := h.profileService.GetUserIDFromToken(c, tokenHeader)
	if err != nil {
		h.logger.Errorf("failed to get user id from token: %v", err)
		switch {
		case errors.Is(err, application.ErrTokenInvalidOrExpired):
			response.NewErrorResponse(c, http.StatusUnauthorized, "Токен недействителен или истёк")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		h.logger.Errorf("Failed to get file from form: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, "Ошибка получения файла")
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get profile: %v", err)
		switch {
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
			return
		case errors.Is(err, application.ErrProfileNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Профиль пользователя не найден")
			return
		case errors.Is(err, application.ErrInvalidUUIDFormat):
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный формат UUID")
			return
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	avatarURL, err := h.profileService.UploadAvatar(c.Request.Context(), userID, file)
	if err != nil {
		h.logger.Errorf("Failed to upload avatar: %v", err)
		switch {
		case errors.Is(err, application.ErrInvalidFileType):
			response.NewErrorResponse(c, http.StatusBadRequest, "Неподдерживаемый тип файла")
		case errors.Is(err, application.ErrFileSizeTooLarge):
			response.NewErrorResponse(c, http.StatusBadRequest, "Размер файла слишком большой")
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	profile.AvatarURL = avatarURL

	updatedProfile, err := h.profileService.UpdateProfile(c.Request.Context(), profile)
	if err != nil {
		h.logger.Errorf("Failed to update profile: %v", err)
		switch {
		case errors.Is(err, application.ErrInvalidAge):
			response.NewErrorResponse(c, http.StatusBadRequest, "Неверный возраст")
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	profileDTO := &dto.Profile{
		UserID:      updatedProfile.UserID,
		Username:    updatedProfile.Username,
		Description: updatedProfile.Description,
		Age:         updatedProfile.Age,
		Location:    updatedProfile.Location,
		AvatarURL:   updatedProfile.AvatarURL,
		CreatedAt:   updatedProfile.CreatedAt,
		UpdatedAt:   updatedProfile.UpdatedAt,
	}

	response.NewSuccessResponse(c, http.StatusOK, "Аватар загружен успешно", profileDTO)
}
