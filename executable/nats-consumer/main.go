package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	cp "github.com/digisan/cli-prompt"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/nats-io/nats.go"
)

var mCfg map[string]interface{}
var err error

// use outter mCfg
func S(name string) string {
	if v, ok := mCfg[name]; ok {
		return v.(string)
	}
	lk.Log("No argument [%s] in config file\n", name)
	return ""
}
func B(name string) bool {
	if v, ok := mCfg[name]; ok {
		return v.(bool)
	}
	lk.Log("No argument [%s] in config file\n", name)
	return false
}
func I(name string) int {
	if v, ok := mCfg[name]; ok {
		return int(v.(float64))
	}
	lk.Log("No argument [%s] in config file\n", name)
	return 0
}

func main() {

	configPtr := flag.String("c", "./config.json", "config(json) file path")
	flag.Parse()

	mCfg, err = cp.PromptConfig(*configPtr)
	lk.FailOnErr("Invalid JSON config file@ [%v]", err)

	if mCfg != nil {
		fmt.Println("Running...")
	}

	// ------------------------------------------ //

	// Connect to NATS
	url := fmt.Sprintf("nats://%s:%d", S("NatsHost"), I("NatsPort"))
	nc, err := nats.Connect(url)
	lk.FailOnErr("%v", err)
	defer nc.Close()

	js, err := nc.JetStream()
	lk.FailOnErr("%v", err)

	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	subject := S("SubStream") + "." + S("SubSubject")
	sub, err := js.PullSubscribe(subject, S("Durable"), nats.PullMaxWaiting(128))
	lk.FailOnErr("%v", err)

	duration, err := time.ParseDuration(S("FetchFrequency"))
	lk.FailOnErr("%v", err)

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
					if err := msg.Ack(); err == nil {
						receive(msg)
					} else {
						lk.Warn("%v", err)
					}
				}
			}
		}()
		wg.Wait()
	}()
	goto AGAIN
}

func chkHashMD5(msg *nats.Msg) bool {
	md5str := msg.Header["FileMD5"][0]
	// fmt.Println("md5str in header:", md5str)
	h := md5.New()
	_, err = io.Copy(h, bytes.NewReader(msg.Data))
	lk.FailOnErr("%v", err)
	return fmt.Sprintf("%x", h.Sum(nil)) == md5str
}

func receive(msg *nats.Msg) {

	///
	// Dump Header
	fmt.Printf("\n-------------------------------------------------------------------\n")
	spew.Dump(msg.Header)
	fmt.Printf("-------------------------------------------------------------------\n\n")
	///

	if !chkHashMD5(msg) {
		lk.Warn("%v", "Hash MD5 is NOT correct")
		return
	}

	format := msg.Header["Format"][0]
	// fmt.Println(format)

	switch format {
	case ".json":
		fmt.Println(string(jt.Fmt(msg.Data, "  ")))
	default:
		fmt.Println(string(msg.Data))
	}
}
