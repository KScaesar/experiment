package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	log.SetFlags(log.Llongfile)
	js := ExampleJetStreamManager()

	_, err := js.Subscribe("bar", func(msg *nats.Msg) {
		meta, _ := msg.Metadata()
		streamSeq := meta.Sequence.Stream
		fmt.Printf("Stream Sequence  : %v\n", streamSeq)
		fmt.Printf("Consumer Sequence: %v\n", meta.Sequence.Consumer)
		fmt.Printf("Consumer Header: %v\n", msg.Header)
		fmt.Printf("msg body: %v\n", string(msg.Data))
		fmt.Println()

		if streamSeq == 2 {
			// https://stackoverflow.com/questions/76715655/nats-try-to-add-a-consumer-attached-to-js-event-advisory-consumer-max-deliverie
			// https://docs.nats.io/using-nats/developer/develop_jetstream/consumers#dead-letter-queues-type-functionality
			msg.Nak()
		} else {
			msg.Ack()
		}
	}, nats.ManualAck(), nats.MaxDeliver(5))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	time.Sleep(2 * time.Second)
}

func ExampleJetStreamManager() nats.JetStreamContext {
	nc, err := nats.Connect("nats://127.0.0.1:4222,nats://127.0.0.1:4223,nats://127.0.0.1:4224", DefaultNatsOption)
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	// defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}

	// https://docs.nats.io/using-nats/developer/develop_jetstream/model_deep_dive
	stream, err := js.StreamNameBySubject("foo")
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	log.Printf("info: stream: %v", stream)

	// https://docs.nats.io/nats-concepts/jetstream/consumers
	_, err = js.AddConsumer("exp_stream", &nats.ConsumerConfig{
		Durable:       "exp_consumer",
		AckPolicy:     nats.AckExplicitPolicy,
		DeliverPolicy: nats.DeliverAllPolicy,
		// MaxDeliver:    1,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}
	// log.Printf("info: consumer: %#v", consumer)

	return js
}

func DefaultNatsOption(options *nats.Options) error {
	options.AllowReconnect = true
	options.User = "devUser"
	options.Password = "123456"
	return nil
}
