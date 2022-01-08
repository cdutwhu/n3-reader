package n3reader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/nats-io/nats.go"
)

var (
	mFn = map[string]func(nats.JetStreamContext, string, *os.File, string) error{
		".json": pubJson,
		".run":  pubBytes,
	}

	pubMethod = func(filetype string) func(nats.JetStreamContext, string, *os.File, string) error {
		filetype = "." + strings.TrimPrefix(filetype, ".")
		if fn, ok := mFn[filetype]; ok {
			return fn
		}
		lk.Warn("<%s> file type has no specific publish method, <pubBytes> applies", filetype)
		return pubBytes
	}
)

func meta2header(meta string) nats.Header {
	m := make(map[string]interface{})
	lk.FailOnErr("%v", json.Unmarshal([]byte(meta), &m))
	h := nats.Header{}
	for k, v := range m {
		h.Add(k, fmt.Sprint(v))
	}
	return h
}

func publish(js nats.JetStreamContext, subj string, header nats.Header, data []byte) error {

	msg := &nats.Msg{
		Subject: subj,
		Header:  header,
		Data:    data,
	}

	ack, err := js.PublishMsg(msg)
	//ack, err := js.Publish(subj, data)

	if err != nil {
		return err
	}
	lk.Log("ACK: %s", ack.Stream)
	return err
}

func pubBytes(js nats.JetStreamContext, subj string, f *os.File, meta string) error {

	header := meta2header(meta)
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	// publish action
	return publish(js, subj, header, data)
}

func pubJson(js nats.JetStreamContext, subj string, f *os.File, meta string) error {

	header := meta2header(meta)

	hint := map[bool]string{
		true:  "Array JSON Fetched",
		false: "Object JSON Fetched",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cOut, arr := jt.ScanObject(ctx, f, false, true, jt.OUT_ORI)
	lk.Log("%s", hint[arr])

	for result := range cOut {
		if result.Err != nil {
			lk.Warn("%v", result.Err)
			return result.Err
		}

		// make publish data
		data := jt.Minimize(fmt.Sprintf(`{"meta":%s, "data":%s}`, meta, result.Obj), true)

		// publish action
		if err := publish(js, subj, header, []byte(data)); err != nil {
			return err
		}
	}

	return nil
}

func (nr *NatsReader) Publish(file, fwMeta string) error {

	lk.Log("Publishing: %s", file)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// merge file watcher meta & nats reader meta
	meta := jt.MergeSgl(fwMeta, nr.exMeta())
	ext := filepath.Ext(file)
	return pubMethod(ext)(nr.js, nr.subject, f, meta)
}
