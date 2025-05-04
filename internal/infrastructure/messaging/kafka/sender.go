package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/notification/entity"
	"gitlab.mai.ru/cicada-chess/backend/user-service/internal/domain/notification/interfaces"
)

type KafkaNotificationSender struct {
	producer Producer
	topic    string
	logger   *logrus.Logger
}

type Producer interface {
	Send(topic string, message []byte) error
	Close() error
}

func NewKafkaNotificationSender(producer Producer, topic string, logger *logrus.Logger) interfaces.NotificationSender {
	return &KafkaNotificationSender{
		producer: producer,
		topic:    topic,
		logger:   logger,
	}
}

func (k *KafkaNotificationSender) SendAccountConfirmation(ctx context.Context, userId, email, token string) error {
	event := entity.NewAccountConfirmationEvent(email, token)

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := k.producer.Send(k.topic, eventJSON); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	k.logger.Printf("Account confirmation sent successfully: %v", email)
	return nil
}

func (k *KafkaNotificationSender) SendPasswordReset(ctx context.Context, userId, email, token string) error {
	event := entity.NewResetPasswordEvent(email, token)

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	if err := k.producer.Send(k.topic, eventJSON); err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	k.logger.Printf("Password reset notification sent successfully: %v", email)

	return nil
}
