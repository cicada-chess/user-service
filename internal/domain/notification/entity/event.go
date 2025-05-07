package entity

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

func NewAccountConfirmationEvent(email, username, token string) *Event {
	return &Event{
		Type:  AccountConfirmation,
		Email: email,
		Payload: map[string]string{
			"link":     generateConfirmationLink(token),
			"username": username,
		},
	}
}

func NewResetPasswordEvent(email, username, token string) *Event {
	return &Event{
		Type:  PasswordReset,
		Email: email,
		Payload: map[string]string{
			"link":     generatePasswordResetLink(token),
			"username": username,
		},
	}
}

func generateConfirmationLink(token string) string {
	return "https://cicada-chess.ru:8081/auth/confirm-account?token=" + token
}

func generatePasswordResetLink(token string) string {
	return "https://cicada-chess.ru:8081/auth/reset-password?token=" + token
}
