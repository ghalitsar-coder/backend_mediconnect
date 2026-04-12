package messaging

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ wraps an AMQP connection and channel.
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// ConnectRabbitMQ establishes a connection to RabbitMQ and opens a channel.
func ConnectRabbitMQ(amqpURL string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	slog.Info("Successfully connected to RabbitMQ")
	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close gracefully closes the RabbitMQ channel and connection.
func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
