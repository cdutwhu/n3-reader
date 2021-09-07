package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	subSubject = "STREAM-1.test"
	pubSubject = "STREAM-1.test-receipt"
)

func main() {

	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	sub, _ := js.PullSubscribe(subSubject, "consumer-test", nats.PullMaxWaiting(128))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msgs, _ := sub.Fetch(100, nats.Context(ctx))
		for _, msg := range msgs {
			msg.Ack()
			data := string(msg.Data)
			fmt.Println(data)
			receipt(js, data)
		}
	}
}

func receipt(js nats.JetStreamContext, data string) {
	_, err := js.Publish(pubSubject, []byte("OK"))
	if err != nil {
		log.Fatal(err)
	}
}
