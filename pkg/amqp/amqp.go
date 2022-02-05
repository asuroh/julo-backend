package amqp

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// IQueue ...
type IQueue interface {
	PushQueue(data map[string]interface{}, types string) error
	PushQueueReconnect(url string, data map[string]interface{}, types, deadLetterKey string) (*amqp.Connection, *amqp.Channel, error)
	PushDLQueueReconnect(url string, data map[string]interface{}, types string) (*amqp.Connection, *amqp.Channel, error)
	PushQueueDelayReconnect(url string, data map[string]interface{}, delayedKey, incomingKey, deadLetterKey string, ttl int64) (*amqp.Connection, *amqp.Channel, error)
}

var (
	// UpdateBalanceExchange ...
	UpdateBalanceExchange = "update_balance.exchange"
	// UpdateBalance ...
	UpdateBalance = "update_balance.incoming.queue"
	// UpdateBalanceDeadLetter ...
	UpdateBalanceDeadLetter = "update_balance.deadletter.queue"
)

// queue ...
type queue struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewQueue ...
func NewQueue(conn *amqp.Connection, channel *amqp.Channel) IQueue {
	return &queue{
		Connection: conn,
		Channel:    channel,
	}
}

// PushQueue ...
func (m queue) PushQueue(data map[string]interface{}, types string) error {
	queue, err := m.Channel.QueueDeclare(types, true, false, false, false, nil)
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = m.Channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	return err
}

// PushQueueReconnect ...
func (m queue) PushQueueReconnect(url string, data map[string]interface{}, types, deadLetterKey string) (*amqp.Connection, *amqp.Channel, error) {
	if m.Connection != nil {
		if m.Connection.IsClosed() {
			c := Connection{
				URL: url,
			}
			newConn, newChannel, err := c.Connect()
			if err != nil {
				return nil, nil, err
			}
			m.Connection = newConn
			m.Channel = newChannel
		}
	} else {
		c := Connection{
			URL: url,
		}
		newConn, newChannel, err := c.Connect()
		if err != nil {
			return nil, nil, err
		}
		m.Connection = newConn
		m.Channel = newChannel
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": deadLetterKey,
	}

	queue, err := m.Channel.QueueDeclare(types, true, false, false, false, args)
	if err != nil {
		return nil, nil, err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, nil, nil
	}

	err = m.Channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	return m.Connection, m.Channel, err
}

// PushDLQueueReconnect ...
func (m queue) PushDLQueueReconnect(url string, data map[string]interface{}, types string) (*amqp.Connection, *amqp.Channel, error) {
	if m.Connection != nil {
		if m.Connection.IsClosed() {
			c := Connection{
				URL: url,
			}
			newConn, newChannel, err := c.Connect()
			if err != nil {
				return nil, nil, err
			}
			m.Connection = newConn
			m.Channel = newChannel
		}
	} else {
		c := Connection{
			URL: url,
		}
		newConn, newChannel, err := c.Connect()
		if err != nil {
			return nil, nil, err
		}
		m.Connection = newConn
		m.Channel = newChannel
	}

	queue, err := m.Channel.QueueDeclare(types, true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, nil, nil
	}

	err = m.Channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	return m.Connection, m.Channel, err
}

// PushQueueDelayReconnect ...
func (m queue) PushQueueDelayReconnect(url string, data map[string]interface{}, delayedKey, incomingKey, deadLetterKey string, ttl int64) (*amqp.Connection, *amqp.Channel, error) {
	if m.Connection != nil {
		if m.Connection.IsClosed() {
			c := Connection{
				URL: url,
			}
			newConn, newChannel, err := c.Connect()
			if err != nil {
				return nil, nil, err
			}
			m.Connection = newConn
			m.Channel = newChannel
		}
	} else {
		c := Connection{
			URL: url,
		}
		newConn, newChannel, err := c.Connect()
		if err != nil {
			return nil, nil, err
		}
		m.Connection = newConn
		m.Channel = newChannel
	}

	m.Channel.QueueDelete(delayedKey, false, false, false)
	args := amqp.Table{
		"x-message-ttl":             ttl,
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": incomingKey,
	}
	queue, err := m.Channel.QueueDeclare(delayedKey, true, false, false, false, args)
	if err != nil {
		return nil, nil, err
	}

	args = amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": deadLetterKey,
	}
	_, err = m.Channel.QueueDeclare(incomingKey, true, false, false, false, args)
	if err != nil {
		return nil, nil, err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}

	err = m.Channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	return m.Connection, m.Channel, err
}
