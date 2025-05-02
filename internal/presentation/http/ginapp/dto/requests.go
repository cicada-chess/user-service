package dto

type GetAllUsersRequest struct {
	Page   string `form:"page"`
	Limit  string `form:"limit"`
	Search string `form:"search"`
	SortBy string `form:"sort_by"`
	Order  string `form:"order"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateRatingRequest struct {
	Delta int `json:"delta" binding:"required"`
}

type UpdateInfoRequest struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Rating   *int    `json:"rating"`
	Role     *int    `json:"role"`
	IsActive *bool   `json:"is_active"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Description *string `json:"description"`
	Age         *int    `json:"age"`
	Location    *string `json:"location"`
	AvatarURL   *string `json:"avatar_url"`
}
