// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// var (
// 	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
// 	exchange     = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
// 	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
// 	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
// 	bindingKey   = flag.String("key", "test-key", "AMQP binding key")
// 	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
// 	lifetime     = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
// )

func init() {
	// flag.Parse()
}

func main() {
	// https://github.com/streadway/amqp/blob/master/_examples/simple-consumer/consumer.go
	conn, err := NewConnection()
	if err != nil {
		panic(err)
	}

	m1, err := NewConsumerManager(conn)
	if err != nil {
		panic(err)
	}
	m2, err := NewConsumerManager(conn)
	if err != nil {
		panic(err)
	}

	go func() {
		for i := 0; i < 5; i++ {
			user := fmt.Sprintf("%v", i)
			queue := fmt.Sprintf("user-%v", user)

			err = m1.Consume("yyy", "fanout", queue, "broadcast") // 1個 queue 有多個消費者, 造型奇怪現象
			if err != nil {
				log.Fatalf("%s", err)
			}

			bindingKey := fmt.Sprintf("login.%v", user)
			err = m2.Consume("xxx", "topic", queue, bindingKey) // 1個 queue 有多個消費者, 造型奇怪現象
			if err != nil {
				log.Fatalf("%s", err)
			}

			log.Printf("\n")
		}
	}()

	time.Sleep(60 * time.Minute)
	// time.Sleep(60 * time.Second)
	if err := m1.ShutdownAll(); err != nil {
		log.Fatalf("ShutdownAll: %v", err)
		return
	}
	if err := m2.ShutdownAll(); err != nil {
		log.Fatalf("ShutdownAll: %v", err)
		return
	}
	if err := conn.Close(); err != nil {
		log.Fatalf("AMQP connection close error: %s", err)
		return
	}

	log.Printf("end!\n")
}

func NewConnection() (*amqp.Connection, error) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	log.Printf("dialing %q", amqpURI)

	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-conn.NotifyClose(make(chan *amqp.Error)))
	}()
	return conn, err
}

func NewConsumerManager(conn *amqp.Connection) (*ConsumerManager, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	manager := &ConsumerManager{
		channel:     channel,
		mu:          sync.RWMutex{},
		consumerAll: make([]string, 0),
	}
	return manager, nil
}

type ConsumerManager struct {
	channel     *amqp.Channel
	mu          sync.RWMutex
	consumerAll []string
}

func (c *ConsumerManager) Consume(exchange, exchangeType, queueName, key string) error {
	if err := c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	consumerName := fmt.Sprintf("%v-%v-%v", exchange, exchangeType, queueName)
	c.consumerAll = append(c.consumerAll, consumerName)

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", consumerName)
	deliveries, err := c.channel.Consume(
		queue.Name,   // name
		consumerName, // consumerTag,
		false,        // noAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, consumerName)
	return nil
}

func (c *ConsumerManager) ShutdownAll() error {
	// will close() the deliveries channel
	for _, consumerName := range c.consumerAll {
		if err := c.channel.Cancel(consumerName, true); err != nil {
			return fmt.Errorf("ConsumerManager cancel failed: %s", err)
		}
	}

	defer log.Printf("AMQP shutdown OK")
	return nil
}

func handle(deliveries <-chan amqp.Delivery, ctag string) {
	for d := range deliveries {
		log.Printf(
			"ctag: [%v], got %dB, DeliveryTag: [%v], payload: %q",
			ctag,
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed: %v", ctag)
}
