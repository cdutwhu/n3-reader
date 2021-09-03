package filewatcher

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
)

type Watcher struct {
	id           string // meta
	name         string // meta
	format       string // meta
	watchFolder  string
	watchFileExt string
	recursive    bool
	inclHidden   bool
	ignore       string
	watcher      *watcher.Watcher
	interval     time.Duration
	Event        IWatchEvent
}

func (w *Watcher) meta(filename string) string {
	return fmt.Sprintf(`{
		"ReaderID": "%s",
		"ReaderName": "%s",
		"SourceFileFormat": "%s",				
		"SourceFileName":"%s",		
		"ReadTimestampUTC":"%s"
	}`, w.id, w.name, w.format, filename, time.Now().UTC().Format(time.RFC3339))
}

func NewFileWatcher(options ...Option) (*Watcher, error) {
	w := &Watcher{Event: &dftEvent{}}
	if err := w.setOption(options...); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Watcher) Close(cleanup func(watched string)) {
	if w.watcher != nil {
		w.watcher.Close() // stop watcher
	}
	if cleanup != nil {
		cleanup(w.watchFolder)
	}
}

func (w *Watcher) start() error {

	go func() {
		for {
			select {
			case event := <-w.watcher.Event:
				if !event.IsDir() {
					var (
						path = event.Path
						meta = w.meta(event.Path)
					)
					switch event.Op {
					case watcher.Remove:
						w.Event.OnDelete(path, meta, time.Now())

					case watcher.Create:
						w.Event.OnCreate(path, meta, event.ModTime())

					case watcher.Write:
						w.Event.OnWrite(path, meta, event.ModTime())
					}
				}

			case err := <-w.watcher.Error:
				w.Event.OnError(err, time.Now())
				return

			case <-w.watcher.Closed:
				w.Event.OnClose(time.Now())
				return
			}
		}
	}()

	// Start the watching process.
	return w.watcher.Start(w.interval)
}

func (w *Watcher) StartWait(cleanup func(watched string)) {

	// signal handler for shutdown
	closed := make(chan struct{})
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nwatcher shutting down")
		w.Close(cleanup) // stop watcher
		fmt.Println("watcher closed")
		close(closed) // release process
	}()

	w.start()
	<-closed
}
