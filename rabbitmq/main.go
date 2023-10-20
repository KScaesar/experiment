// This example declares a durable Exchange, an ephemeral (auto-delete) Queue,
// binds the Queue to the Exchange with a binding key, and consumes every
// message published to that Exchange with that routing key.
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

	manager := NewManager()
	var channelAll []*amqp.Channel

	go func() {
		for i := 0; i < 3; i++ {
			user := fmt.Sprintf("user%v", i)
			queueName := fmt.Sprintf("notify.%v", user)

			ex1 := "broadcast"
			ex1Type := "fanout"
			key1 := "broadcast.lv0.*"
			_, channel1, err1 := NewQueueByConnection(conn, ex1, ex1Type, queueName, key1)
			if err1 != nil {
				log.Fatalf("%s", err)
			}

			ex2 := "condition"
			ex2Type := "topic"
			key2 := "*.*.*"
			_, channel2, err2 := NewQueueByConnection(conn, ex2, ex2Type, queueName, key2)
			if err2 != nil {
				log.Fatalf("%s", err2)
			}

			ex3 := "single"
			ex3Type := "direct"
			key3 := user
			_, channel3, err3 := NewQueueByConnection(conn, ex3, ex3Type, queueName, key3)
			if err3 != nil {
				log.Fatalf("%s", err3)
			}

			channelAll = append(channelAll, []*amqp.Channel{channel1, channel2, channel3}...)
		}

		for i := 0; i < 3; i++ {
			user := fmt.Sprintf("user%v", i)
			queueName := fmt.Sprintf("notify.%v", user)

			builder := strings.Builder{}
			builder.WriteString(queueName)
			builder.WriteString("-")
			builder.WriteString("worker")
			consumerName := builder.String()
			consumers, err := NewConsumerAllBySingleChannel(channelAll[i], queueName, consumerName, 1)
			if err != nil {
				log.Fatalf("%v", err)
				return
			}
			manager.AddConsumerAndServeConsume(user, consumers...)
		}
	}()

	time.Sleep(60 * time.Minute)

	manager.StopConsumerAll()

	if err := conn.Close(); err != nil {
		log.Fatalf("AMQP connection close error: %s", err)
		return
	}

	log.Printf("end!\n")
}

func NewManager() *Manager {
	return &Manager{}
}

type Manager struct {
	consumers sync.Map // owner:[]Consumer
}

func (m *Manager) AddConsumerAndServeConsume(owner string, consumers ...Consumer) {
	m.consumers.Store(owner, consumers)
	for _, consumer := range consumers {
		consumer := consumer
		go func() {
			consumer.ServeConsume()
		}()
	}
}

func (m *Manager) StopConsumerByOwner(owner string) {
	value, exist := m.consumers.LoadAndDelete(owner)
	if !exist {
		return
	}

	for _, consumer := range value.([]Consumer) {
		err := consumer.Shutdown()
		if err != nil {
			log.Printf("consumer=%v: shutdown: %v", consumer.name, err)
		}
	}
	return
}

func (m *Manager) StopConsumerAll() {
	wg := sync.WaitGroup{}
	m.consumers.Range(func(key, value any) bool {
		for _, consumer := range value.([]Consumer) {
			consumer := consumer
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := consumer.Shutdown()
				if err != nil {
					log.Printf("consumer=%v: shutdown: %v", consumer.name, err)
				}
			}()
		}
		return true
	})
	wg.Wait()
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

func NewQueueByChannel(ch *amqp.Channel, exchange, exchangeType, queueName, key string) (amqp.Queue, error) {
	if err := ch.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return amqp.Queue{}, fmt.Errorf("Exchange Declare: %s", err)
	}

	queue, err := ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("Queue Declare: %s", err)
	}

	if err = ch.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return amqp.Queue{}, fmt.Errorf("Queue Bind: %s", err)
	}

	return queue, err
}

func NewQueueByConnection(conn *amqp.Connection, exchange, exchangeType, queueName, key string) (amqp.Queue, *amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return amqp.Queue{}, nil, fmt.Errorf("NewChannel: %s", err)
	}

	queue, err := NewQueueByChannel(channel, exchange, exchangeType, queueName, key)
	if err != nil {
		return amqp.Queue{}, nil, fmt.Errorf("NewQueue: %s", err)
	}

	return queue, channel, err
}

func NewConsumerAllBySingleChannel(ch *amqp.Channel, queueName string, consumerName string, consumerQty int) ([]Consumer, error) {
	var consumers []Consumer
	for i := 0; i < consumerQty; i++ {
		cTag := consumerName + strconv.Itoa(i)

		log.Printf("starting Consume (consumer tag %q)", cTag)
		deliveries, err := ch.Consume(
			queueName, // name
			cTag,      // consumerTag,
			false,     // noAck
			false,     // exclusive
			false,     // noLocal
			false,     // noWait
			nil,       // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("NewConsumerAllBySingleChannel: %s", err)
		}

		consumers = append(consumers, Consumer{
			channel:  ch,
			delivery: deliveries,
			name:     cTag,
			done:     make(chan struct{}),
		})
	}

	return consumers, nil
}

type Consumer struct {
	channel  *amqp.Channel        // for shutdown
	delivery <-chan amqp.Delivery // for consume
	name     string               // for consume
	done     chan struct{}
}

func (c *Consumer) ServeConsume() {
	for d := range c.delivery {
		log.Printf(
			"ctag: [%v], got %dB, DeliveryTag: [%v], payload: %q, key: %v",
			c.name,
			len(d.Body),
			d.DeliveryTag,
			d.Body,
			d.RoutingKey,
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed: %v", c.name)
	close(c.done)
}

func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.name, true); err != nil {
		return fmt.Errorf("concumer=%v: cancel: %v", c.name, err)
	}
	<-c.done
	return nil
}
