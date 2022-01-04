package filewatcher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type IWatchEvent interface {
	OnCreate(path, fwMeta string, t time.Time) error
	OnWrite(path, fwMeta string, t time.Time) error
	OnDelete(path, fwMeta string, t time.Time) error
	OnError(err error, t time.Time) error
	OnClose(t time.Time) error
}

// default event example

type Event struct{}

func (evt *Event) OnCreate(path, fwMeta string, t time.Time) error {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(fwMeta), &m)
	fmt.Printf("\nfile: %s\n[Created]: %s\n", path, t)
	spew.Dump(m)
	return nil
}

func (evt *Event) OnWrite(path, fwMeta string, t time.Time) error {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(fwMeta), &m)
	fmt.Printf("\nfile: %s\n[Modified]: %s\n", path, t)
	spew.Dump(m)
	return nil
}

func (evt *Event) OnDelete(path, fwMeta string, t time.Time) error {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(fwMeta), &m)
	fmt.Printf("\nfile: %s\n[Deleted]: %s\n", path, t)
	spew.Dump(m)
	return nil
}

func (evt *Event) OnError(err error, t time.Time) error {
	fmt.Println("\tFile-Watcher error occurred: ", err, t)
	return nil
}

func (evt *Event) OnClose(t time.Time) error {
	fmt.Println("\tFile-Watcher closed at ", t)
	return nil
}
