package kafka

import (
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer  sarama.SyncProducer
	brokers   []string
	mutex     sync.Mutex
	isRunning bool
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &KafkaProducer{
		producer:  producer,
		brokers:   brokers,
		isRunning: true,
	}, nil
}

func (p *KafkaProducer) Send(topic string, message []byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.isRunning {
		return fmt.Errorf("producer is closed")
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Message sent to topic: %s, partition: %d, offset: %d", topic, partition, offset)
	return nil
}

func (p *KafkaProducer) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.isRunning {
		return nil
	}

	p.isRunning = false
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}

	log.Println("Kafka producer closed")
	return nil
}
