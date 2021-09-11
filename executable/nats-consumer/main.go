package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	cp "github.com/digisan/cli-prompt"
	jt "github.com/digisan/json-tool"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

var mc map[string]interface{}
var err error

// use outter mc
func S(name string) string {
	return mc[name].(string)
}
func I(name string) int {
	return int(mc[name].(float64))
}

func main() {

	configPtr := flag.String("c", "./config.json", "config(json) file path")
	flag.Parse()

	mc, err = cp.PromptConfig(*configPtr)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Invalid config file as JSON format"))
	}
	if mc != nil {
		fmt.Println("Running...")
	}

	// ------------------------------------------ //

	// Connect to NATS
	url := fmt.Sprintf("nats://%s:%d", S("NatsHost"), I("NatsPort"))
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	

	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	subject := S("subStream") + "." + S("subSubject")
	sub, err := js.PullSubscribe(subject, S("durable"), nats.PullMaxWaiting(128))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msgs, _ := sub.Fetch(100, nats.Context(ctx))
			for _, msg := range msgs {
				msg.Ack()
				receive(msg.Data)
			}
		}
	}()

	js.Subscribe(subject, func(msg *nats.Msg) {
		msg.Ack()
		receive(msg.Data)
	})

	wg.Wait()
}

func receive(data []byte) {
	fmt.Println(string(jt.Fmt(data, "  ")))
}
