package n3reader

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type n3ReaderEvent struct {
	*Nats4Reader
}

func NewN3ReaderEvent(n4r *Nats4Reader) *n3ReaderEvent {
	return &n3ReaderEvent{n4r}
}

func (event *n3ReaderEvent) OnCreate(path, meta string, t time.Time) {
	event.PubAsJSON(path, meta)
}

func (event *n3ReaderEvent) OnWrite(path, meta string, t time.Time) {
	// event.PubAsJSON(path, meta)
}

func (event *n3ReaderEvent) OnDelete(path, meta string, t time.Time) {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(meta), &m)
	fmt.Printf("\nfile: %s\n[Deleted]: %s\n", path, t)
	spew.Dump(m)
}

func (event *n3ReaderEvent) OnError(err error, t time.Time) {

}

func (event *n3ReaderEvent) OnClose(t time.Time) {
	event.nc.Close()
}
