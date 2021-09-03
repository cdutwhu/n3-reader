package main

import (
	"fmt"
	"os"

	. "github.com/cdutwhu/n3-reader"
	fw "github.com/cdutwhu/n3-reader/file-watcher"
)

func main() {

	cleanup := func(folder string) {
		if err := os.RemoveAll(folder); err != nil {
			panic(err)
		}
		fmt.Printf("%s is removed\n", folder)
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
		// freader.StartWait(cleanup)
	}

	{
		if n4r, err := NewNats4Reader(); err == nil {
			opts := []fw.Option{
				fw.OptID(""),
				fw.OptFormat("json"),
				fw.OptName(""),
				fw.OptWatcher("", "json", "100ms", false, false, ""),
			}
			if freader, err := fw.NewFileWatcher(opts...); err == nil {
				freader.Event = NewN3ReaderEvent(n4r)
				freader.StartWait(cleanup)
			}
		}
	}
}
