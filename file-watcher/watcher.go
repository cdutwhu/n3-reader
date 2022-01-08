package filewatcher

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/radovskyb/watcher"
)

type Watcher struct {
	id         string           // meta
	name       string           // meta
	fileKind   EmFileKind       // meta
	fileType   EmFileType       // meta
	format     []string         // no meta
	folder     string           // no meta
	fileExt    string           // meta
	recursive  bool             // no meta
	inclHidden bool             // no meta
	ignore     string           // no meta
	autodel    bool             // no meta
	watcher    *watcher.Watcher // no meta
	interval   time.Duration    // no meta
	Event      IWatchEvent      // no meta
}

func (w *Watcher) Id() string       { return w.id }
func (w *Watcher) Name() string     { return w.name }
func (w *Watcher) Format() []string { return w.format }
func (w *Watcher) Folder() string   { return w.folder }
func (w *Watcher) FileExt() string  { return w.fileExt }

func (w *Watcher) meta(file string, filekind EmFileKind) string {
	w.fileType = getFileType(file)
	m := map[string]interface{}{
		"ReaderID":         w.id,
		"ReaderName":       w.name,
		"FileKind":         w.fileKind.String(),
		"FileType":         w.fileType.String(),
		"Format":           filepath.Ext(file),
		"Source":           filepath.Base(file),
		"ReadTimestampUTC": time.Now().UTC().Format(time.RFC3339),
	}
	data, err := json.Marshal(m)
	lk.FailOnErr("%v", err)
	return string(data)
}

func NewFileWatcher(options ...Option) (*Watcher, error) {
	w := &Watcher{Event: &Event{}}
	if err := w.setOption(options...); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Watcher) Init(prepare func(watcher *Watcher)) {
	if prepare != nil {
		prepare(w)
	}
}

func (w *Watcher) Close(cleanup func(watcher *Watcher)) {
	if w.watcher != nil {
		w.watcher.Close() // stop watcher
	}
	if cleanup != nil {
		cleanup(w)
	}
}

func (w *Watcher) start() error {

	go func() {
		adPath := ""
		for {
			select {
			case event := <-w.watcher.Event:
				if !event.IsDir() {
					var (
						path = event.Path
						meta = w.meta(event.Path, w.fileKind)
					)
					if HasAnySuffix(path, w.format...) { // only interested in specific format
						switch event.Op {
						case watcher.Remove:
							if adPath != path {
								e := w.Event.OnDelete(path, meta, time.Now())
								lk.WarnOnErr("<OnDelete> Error@ %v", e)
							}

						case watcher.Create:
							e := w.Event.OnCreate(path, meta, event.ModTime())
							lk.WarnOnErr("<OnCreate> Error@ %v", e)
							if w.autodel {
								adPath = path
								lk.FailOnErr("<OnCreate-AutoDelete> Error@ %v", os.Remove(path))
							}

						case watcher.Write:
							e := w.Event.OnWrite(path, meta, event.ModTime())
							lk.WarnOnErr("<OnWrite> Error@ %v", e)
						}
					} else {
						lk.WarnOnErr("<%s> type file is ignored\n", filepath.Ext(path))
					}
				}

			case err := <-w.watcher.Error:
				e := w.Event.OnError(err, time.Now())
				lk.WarnOnErr("<OnError> Error@ %v", e)
				return

			case <-w.watcher.Closed:
				e := w.Event.OnClose(time.Now())
				lk.WarnOnErr("<OnClose> Error@ %v", e)
				return
			}
		}
	}()

	// Start the watching process.
	return w.watcher.Start(w.interval)
}

func (w *Watcher) StartWait(prepare, cleanup func(watcher *Watcher)) {

	lk.Log("watcher is running...")

	w.Init(prepare) // do some preparation

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
