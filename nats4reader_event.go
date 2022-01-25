package n3reader

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type ReaderEvent struct {
	*NatsReader
}

func NewReaderEvent(nr *NatsReader) *ReaderEvent {
	return &ReaderEvent{nr}
}

func (evt *ReaderEvent) OnCreate(path, fwMeta string, t time.Time) error {

	///
	// TODO: split file to parts if file size is bigger
	///

	return evt.Publish(path, fwMeta)
}

func (evt *ReaderEvent) OnWrite(path, fwMeta string, t time.Time) error {
	// evt.PubAsJSON(path, meta)
	return nil
}

func (evt *ReaderEvent) OnDelete(path, fwMeta string, t time.Time) error {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(fwMeta), &m)
	fmt.Printf("\nfile: %s\n[Deleted]: %s\n", path, t)
	spew.Dump(m)
	return nil
}

func (evt *ReaderEvent) OnError(err error, t time.Time) error {
	return err
}

func (evt *ReaderEvent) OnClose(t time.Time) error {
	evt.nc.Close()
	return nil
}
