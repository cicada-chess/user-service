package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
// @Success 200 {object} response.SuccessResponse{data=dto.Profile} "Профиль пользователя"
// @Failure 401 {object} response.ErrorResponse "Ошибка авторизации"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
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

	// Формируем полный URL для аватара, если он есть
	avatarURL := ""
	if profile.AvatarPath != "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		avatarURL = fmt.Sprintf("%s://%s/uploads/avatars/%s", scheme, c.Request.Host, filepath.Base(profile.AvatarPath))
	}

	// Преобразуем в DTO
	profileDTO := &dto.Profile{
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
// @Param request body dto.UpdateProfileRequest true "Данные профиля"
// @Success 200 {object} response.SuccessResponse{data=dto.Profile} "Профиль обновлен"
// @Failure 400 {object} response.ErrorResponse "Неверные данные профиля"
// @Failure 401 {object} response.ErrorResponse "Ошибка авторизации"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /profile [put]
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

	// Преобразуем из DTO в доменную модель
	profile := &entity.Profile{
		UserID:      userID,
		Description: request.Description,
		Age:         request.Age,
		Location:    request.Location,
	}

	updatedProfile, err := h.profileService.CreateOrUpdateProfile(c.Request.Context(), profile)
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

	// Формируем полный URL для аватара, если он есть
	avatarURL := ""
	if updatedProfile.AvatarPath != "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		avatarURL = fmt.Sprintf("%s://%s/%s", scheme, c.Request.Host, updatedProfile.AvatarPath)
	}

	// Преобразуем в DTO для ответа
	profileDTO := &dto.Profile{
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
// @Param description formData string false "Описание профиля"
// @Param age formData integer false "Возраст"
// @Param location formData string false "Местоположение"
// @Success 200 {object} response.SuccessResponse{data=dto.Profile} "Аватар загружен"
// @Failure 400 {object} response.ErrorResponse "Ошибка загрузки файла"
// @Failure 401 {object} response.ErrorResponse "Ошибка авторизации"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
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

	// Получаем файл аватара
	file, err := c.FormFile("avatar")
	if err != nil {
		h.logger.Errorf("Failed to get file from form: %v", err)
		response.NewErrorResponse(c, http.StatusBadRequest, "Ошибка получения файла")
		return
	}

	// Загружаем аватар
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

	// Обновляем другие поля профиля, если они были предоставлены
	profile := &entity.Profile{
		UserID:      userID,
		Description: c.PostForm("description"),
		AvatarPath:  avatarPath,
	}

	// Обрабатываем age
	if ageStr := c.PostForm("age"); ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err == nil {
			profile.Age = age
		}
	}

	profile.Location = c.PostForm("location")

	// Обновляем профиль
	updatedProfile, err := h.profileService.CreateOrUpdateProfile(c.Request.Context(), profile)
	if err != nil {
		h.logger.Errorf("Failed to update profile during avatar upload: %v", err)
		// Продолжаем, так как аватар уже загружен
	}

	// Формируем полный URL для аватара
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	avatarURL := fmt.Sprintf("%s://%s/%s", scheme, c.Request.Host, avatarPath)

	// Преобразуем в DTO для ответа
	profileDTO := &dto.Profile{
		Description: updatedProfile.Description,
		Age:         updatedProfile.Age,
		Location:    updatedProfile.Location,
		AvatarURL:   avatarURL,
		CreatedAt:   updatedProfile.CreatedAt,
		UpdatedAt:   updatedProfile.UpdatedAt,
	}

	response.NewSuccessResponse(c, http.StatusOK, "Аватар загружен успешно", profileDTO)
}
