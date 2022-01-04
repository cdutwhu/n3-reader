package n3reader

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	jt "github.com/digisan/json-tool"
	"github.com/nats-io/nats.go"
)

var (
	mFn = map[string]func(nats.JetStreamContext, string, *os.File, string) error{
		".json": pubJson,
		".run":  pubBytes,
	}
)

func publish(js nats.JetStreamContext, subj string, data []byte) error {
	ack, err := js.Publish(subj, data)
	if err != nil {
		return err
	}
	log.Println("ACK:", ack.Stream)
	return err
}

func pubBytes(js nats.JetStreamContext, subj string, f *os.File, meta string) error {
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	// publish action
	return publish(js, subj, data)
}

func pubJson(js nats.JetStreamContext, subj string, f *os.File, meta string) error {

	msg := map[bool]string{
		true:  "Array JSON Fetched",
		false: "Object JSON Fetched",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cOut, arr := jt.ScanObject(ctx, f, false, true, jt.OUT_ORI)
	log.Println(msg[arr])

	for result := range cOut {
		if result.Err != nil {
			log.Println(result.Err)
			return result.Err
		}

		// make publish data
		data := jt.Minimize(fmt.Sprintf(`{"meta":%s, "data":%s}`, meta, result.Obj), true)

		// publish action
		if err := publish(js, subj, []byte(data)); err != nil {
			return err
		}
	}

	return nil
}

func (nr *NatsReader) Publish(file, fwMeta string) error {

	log.Println("Publishing:", file)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// merge file watcher meta & nats reader meta
	meta := jt.MergeSgl(fwMeta, nr.exMeta())
	ext := filepath.Ext(file)

	if fn, ok := mFn[ext]; ok {
		return fn(nr.js, nr.subject, f, meta)
	}

	return fmt.Errorf("<%s> file type is NOT supported", ext)
}
