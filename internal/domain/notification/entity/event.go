package entity

// EventType тип события
type EventType string

const (
	AccountConfirmation EventType = "user_registration"
	PasswordReset       EventType = "password_reset_request"
)

type Event struct {
	Type    EventType
	Email   string
	Payload map[string]string
}

func NewAccountConfirmationEvent(email, token string) *Event {
	return &Event{
		Type:  AccountConfirmation,
		Email: email,
		Payload: map[string]string{
			"link": generateConfirmationLink(token),
		},
	}
}

func NewResetPasswordEvent(email, token string) *Event {
	return &Event{
		Type:  PasswordReset,
		Email: email,
		Payload: map[string]string{
			"link": generatePasswordResetLink(token),
		},
	}
}

func generateConfirmationLink(token string) string {
	// НЕ РАБОТАЕТ НАДО ПЕРЕПИСАТЬ
	return "http://localhost:8080/confirm?token=" + token
}

func generatePasswordResetLink(token string) string {
	// НЕ РАБОТАЕТ НАДО ПЕРЕПИСАТЬ
	return "https://localhost:8080/reset-password?token=" + token
}
