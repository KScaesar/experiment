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
	comsumers := make([]*Consumer, 0)
	mu := sync.Mutex{}
	go func() {
		i := 0
		for i < 5 {
			exchange := "user"
			queue := fmt.Sprintf("user-%v", i)
			// bindingKey := fmt.Sprintf("login.%v", queue)
			bindingKey := "login.#"
			consumerTag := fmt.Sprintf("consumer-%v", i)
			c, err := NewConsumer(exchange, queue, bindingKey, consumerTag)
			if err != nil {
				log.Fatalf("%s", err)
			}
			log.Printf("\n")

			mu.Lock()
			comsumers = append(comsumers, c)
			mu.Unlock()
			i++
			// time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(60 * time.Minute)
	// time.Sleep(60 * time.Second)
	mu.Lock()
	for i, c := range comsumers {
		log.Printf("shutting down: consumer-%v\n", i)
		if err := c.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
	}
	mu.Unlock()
	log.Printf("end!\n")
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(exchange, queueName, key, ctag string) (*Consumer, error) {
	amqpURI := "amqp://guest:guest@localhost:5672/"
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange, // name of the exchange
		"topic",  // type
		true,     // durable
		false,    // delete when complete
		false,    // internal
		false,    // noWait
		nil,      // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go handle(deliveries, c.done, ctag)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error, ctag string) {
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
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
