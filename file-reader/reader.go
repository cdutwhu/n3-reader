package filereader

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
)

type Reader struct {
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
	Event        IReaderEvent
}

func (fr *Reader) meta(filename string) string {
	return fmt.Sprintf(`{
		"ReaderID": "%s",
		"ReaderName": "%s",
		"SourceFileFormat": "%s",				
		"SourceFileName":"%s",		
		"ReadTimestampUTC":"%s"
	}`, fr.id, fr.name, fr.format, filename, time.Now().UTC().Format(time.RFC3339))
}

func NewFileReader(options ...Option) (*Reader, error) {
	fr := &Reader{Event: &dftEvent{}}
	if err := fr.setOption(options...); err != nil {
		return nil, err
	}
	return fr, nil
}

func (fr *Reader) Close(cleanup func(watched string)) {
	fr.watcher.Close() // stop watcher
	cleanup(fr.watchFolder)
}

func (fr *Reader) start() error {

	go func() {
		for {
			select {
			case event := <-fr.watcher.Event:
				if !event.IsDir() {
					var (
						path = event.Path
						meta = fr.meta(event.Path)
					)
					switch event.Op {
					case watcher.Remove:
						fr.Event.OnDelete(path, meta, time.Now())

					case watcher.Create:
						fr.Event.OnCreate(path, meta, event.ModTime())

					case watcher.Write:
						fr.Event.OnWrite(path, meta, event.ModTime())
					}
				}

			case err := <-fr.watcher.Error:
				fr.Event.OnError(err, time.Now())
				return

			case <-fr.watcher.Closed:
				fr.Event.OnClose(time.Now())
				return
			}
		}
	}()

	// Start the watching process.
	return fr.watcher.Start(fr.interval)
}

func (fr *Reader) StartWait(cleanup func(watched string)) {

	// signal handler for shutdown
	closed := make(chan struct{})
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nreader shutting down")
		fr.Close(cleanup) // stop watcher
		fmt.Println("reader closed")
		close(closed) // release process
	}()

	fr.start()
	<-closed
}
