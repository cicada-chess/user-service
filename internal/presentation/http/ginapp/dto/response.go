package dto

type SuccessResponseWithoutData struct {
	Success string `json:"success"`
	Message string `json:"message,omitempty"`
}
