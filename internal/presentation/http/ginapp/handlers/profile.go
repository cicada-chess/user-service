package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "gitlab.mai.ru/cicada-chess/backend/user-service/docs"
	application "gitlab.mai.ru/cicada-chess/backend/user-service/internal/application/profile"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/profile/entity"
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
		h.logger.Error("Authorization header is missing")
		response.NewErrorResponse(c, http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	userID, err := h.profileService.GetUserIDFromToken(c, tokenHeader)
	if err != nil {
		h.logger.Errorf("Failed to get user ID from token: %v", err)
		response.NewErrorResponse(c, http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get profile: %v", err)
		switch {
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, "Ошибка получения профиля")
		}
		return
	}

	avatarURL := ""
	if profile.AvatarPath != "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		avatarURL = fmt.Sprintf("%s://%s/uploads/avatars/%s", scheme, c.Request.Host, filepath.Base(profile.AvatarPath))
	}

	profileDTO := &dto.Profile{
		UserID:      profile.UserID,
		Description: profile.Description,
		Age:         profile.Age,
		Location:    profile.Location,
		AvatarURL:   avatarURL,
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
		h.logger.Error("Authorization header is missing")
		response.NewErrorResponse(c, http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	userID, err := h.profileService.GetUserIDFromToken(c, tokenHeader)
	if err != nil {
		h.logger.Errorf("Failed to get user ID from token: %v", err)
		response.NewErrorResponse(c, http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	var request dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Errorf("Failed to bind profile update request: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, "Неверные данные запроса")
		return
	}

	userProfile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user profile: %v", err)
		switch {
		case errors.Is(err, application.ErrUserNotFound):
			response.NewErrorResponse(c, http.StatusNotFound, "Пользователь не найден")
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, "Ошибка получения профиля")
		}
		return
	}

	if request.Description != nil {
		userProfile.Description = *request.Description
	}
	if request.Age != nil {
		userProfile.Age = *request.Age
	}
	if request.Location != nil {
		userProfile.Location = *request.Location
	}
	profile := &entity.Profile{
		UserID:      userID,
		Description: userProfile.Description,
		Age:         userProfile.Age,
		Location:    userProfile.Location,
		AvatarPath:  userProfile.AvatarPath,
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
			response.NewErrorResponse(c, http.StatusInternalServerError, "Ошибка обновления профиля")
		}
		return
	}

	avatarURL := ""
	if updatedProfile.AvatarPath != "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		avatarURL = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, updatedProfile.AvatarPath)
	}

	profileDTO := &dto.Profile{
		UserID:      updatedProfile.UserID,
		Description: updatedProfile.Description,
		Age:         updatedProfile.Age,
		Location:    updatedProfile.Location,
		AvatarURL:   avatarURL,
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
		h.logger.Error("Authorization header is missing")
		response.NewErrorResponse(c, http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	userID, err := h.profileService.GetUserIDFromToken(c, tokenHeader)
	if err != nil {
		h.logger.Errorf("Failed to get user ID from token: %v", err)
		response.NewErrorResponse(c, http.StatusUnauthorized, "Ошибка авторизации")
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		h.logger.Errorf("Failed to get file from form: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, "Ошибка получения файла")
		return
	}

	avatarPath, err := h.profileService.UploadAvatar(c.Request.Context(), userID, file)
	if err != nil {
		h.logger.Errorf("Failed to upload avatar: %v", err)
		switch {
		case errors.Is(err, application.ErrInvalidFileType):
			response.NewErrorResponse(c, http.StatusBadRequest, "Неподдерживаемый тип файла")
		case errors.Is(err, application.ErrFileSizeTooLarge):
			response.NewErrorResponse(c, http.StatusBadRequest, "Размер файла слишком большой")
		default:
			response.NewErrorResponse(c, http.StatusInternalServerError, "Ошибка загрузки аватара")
		}
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get profile: %v", err)
		response.NewErrorResponse(c, http.StatusInternalServerError, "Ошибка получения профиля")
		return
	}
	profile.AvatarPath = avatarPath

	updatedProfile, err := h.profileService.UpdateProfile(c.Request.Context(), profile)
	if err != nil {
		h.logger.Errorf("Failed to update profile during avatar upload: %v", err)
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	avatarURL := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, avatarPath)

	profileDTO := &dto.Profile{
		UserID:      updatedProfile.UserID,
		Description: updatedProfile.Description,
		Age:         updatedProfile.Age,
		Location:    updatedProfile.Location,
		AvatarURL:   avatarURL,
		CreatedAt:   updatedProfile.CreatedAt,
		UpdatedAt:   updatedProfile.UpdatedAt,
	}

	response.NewSuccessResponse(c, http.StatusOK, "Аватар загружен успешно", profileDTO)
}
