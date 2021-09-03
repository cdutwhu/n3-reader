package filereader

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type IReaderEvent interface {
	OnCreate(path, meta string, t time.Time)
	OnWrite(path, meta string, t time.Time)
	OnDelete(path, meta string, t time.Time)
	OnError(err error, t time.Time)
	OnClose(t time.Time)
}

type dftEvent struct{}

func (e *dftEvent) OnCreate(path, meta string, t time.Time) {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(meta), &m)
	fmt.Printf("\nfile: %s\n[Created]: %s\n", path, t)
	spew.Dump(m)
}

func (e *dftEvent) OnWrite(path, meta string, t time.Time) {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(meta), &m)
	fmt.Printf("\nfile: %s\n[Modified]: %s\n", path, t)
	spew.Dump(m)
}

func (e *dftEvent) OnDelete(path, meta string, t time.Time) {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(meta), &m)
	fmt.Printf("\nfile: %s\n[Deleted]: %s\n", path, t)
	spew.Dump(m)
}

func (e *dftEvent) OnError(err error, t time.Time) {
	fmt.Println("\tFile-Watcher error occurred: ", err, t)
}

func (e *dftEvent) OnClose(t time.Time) {
	fmt.Println("\tFile-Watcher closed at ", t)
	os.RemoveAll("./watched")
}
