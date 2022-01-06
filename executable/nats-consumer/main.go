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
	if v, ok := mc[name]; ok {
		return v.(string)
	}
	log.Fatalf("No argument [%s] in config file\n", name)
	return ""
}
func B(name string) bool {
	if v, ok := mc[name]; ok {
		return v.(bool)
	}
	log.Fatalf("No argument [%s] in config file\n", name)
	return false
}
func I(name string) int {
	if v, ok := mc[name]; ok {
		return int(v.(float64))
	}
	log.Fatalf("No argument [%s] in config file\n", name)
	return 0
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
	subject := S("SubStream") + "." + S("SubSubject")
	sub, err := js.PullSubscribe(subject, S("Durable"), nats.PullMaxWaiting(128))
	if err != nil {
		log.Fatal(err)
	}

	duration, err := time.ParseDuration(S("FetchFrequency"))
	if err != nil {
		log.Fatal(err)
	}
AGAIN:
	func() {
		ctx, cancel := context.WithTimeout(context.Background(), duration)
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
				msgs, _ := sub.Fetch(128, nats.Context(ctx))
				for _, msg := range msgs {
					msg.Ack()
					receive(msg)
				}
			}
		}()
		wg.Wait()
	}()
	goto AGAIN
}

func receive(msg *nats.Msg) {
	switch msg.Header["Format"][0] {
	case "json":
		fmt.Println(string(jt.Fmt(msg.Data, "  ")))
	default:
		fmt.Println(string(msg.Data))
	}
}
