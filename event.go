package n3reader

import (
	"fmt"
	"time"
)

type IFileReaderEvent interface {
	OnCreateWrite(path, meta string, t time.Time)
	OnDelete(path, meta string, t time.Time)
	OnError(err error, t time.Time)
	OnClose(t time.Time)
}

type dftEvent struct{}

func (e *dftEvent) OnCreateWrite(path, meta string, t time.Time) {
	fmt.Printf("\nfile: %s\nmodified: %s\n", path, t)
}

func (e *dftEvent) OnDelete(path, meta string, t time.Time) {
	fmt.Printf("\nfile: %s\ndeleted: %s\n", path, t)
}

func (e *dftEvent) OnError(err error, t time.Time) {
	fmt.Println("\tFile-Watcher error occurred: ", err, t)
}

func (e *dftEvent) OnClose(t time.Time) {
	fmt.Println("\tFile-Watcher closed at ", t)
}
