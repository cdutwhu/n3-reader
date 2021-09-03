package main

import (
	"fmt"
	"os"

	. "github.com/cdutwhu/n3-reader"
	fw "github.com/cdutwhu/n3-reader/file-watcher"
)

func main() {

	prepare := func(w *fw.Watcher) {}
	cleanup := func(w *fw.Watcher) {
		if err := os.RemoveAll(w.Folder); err != nil {
			panic(err)
		}
		fmt.Printf("%s is removed\n", w.Folder)
	}

	{
		// opts := []fw.Option{
		// 	fw.OptID(""),
		// 	fw.OptFormat("json"),
		// 	fw.OptName(""),
		// 	fw.OptWatcher("", "json", "100ms", false, false, ""),
		// }
		// freader, err := fw.NewFileWatcher(opts...)
		// if err != nil {
		// 	panic(err)
		// }
		// freader.StartWait(prepare, cleanup)
	}

	{
		opts := []fw.Option{
			fw.OptID(""),
			fw.OptFormat("json"),
			fw.OptName(""),
			fw.OptWatcher("", "json", "100ms", false, false, ""),
		}
		if freader, err := fw.NewFileWatcher(opts...); err == nil {
			
			opts := []Option{
				OptNatsHostName(""),
				OptNatsPort(0),
				OptNatsClusterName(""),
				OptTopic("fromN3Reader"),
				OptConcurrentFiles(0),
			}
			n3r, err := NewNats4Reader(opts...)
			if err == nil {

				n3r.InitStanConn(freader.Name)
				freader.Event = NewN3ReaderEvent(n3r)
				freader.StartWait(prepare, cleanup)

			} else {
				panic(err)
			}

		} else {
			panic(err)
		}
	}
}
