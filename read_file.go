package n3reader

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"
)

type FileReader struct {
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
	Event        IFileReaderEvent
}

func (fr *FileReader) meta(filename string) string {
	return fmt.Sprintf(`{
		"ReaderID": "%s",
		"ReaderName": "%s",
		"SourceFileFormat": "%s",				
		"SourceFileName":"%s",		
		"ReadTimestampUTC":"%s"
	}`, fr.id, fr.name, fr.format, filename, time.Now().UTC().Format(time.RFC3339))
}

func NewFileReader(options ...Option) (*FileReader, error) {
	fr := &FileReader{Event: &dftEvent{}}
	if err := fr.setOption(options...); err != nil {
		return nil, err
	}
	return fr, nil
}

func (fr *FileReader) Close() {
	fr.watcher.Close()
}

func (fr *FileReader) Start() error {

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
