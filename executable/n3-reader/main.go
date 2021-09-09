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
		if err := os.RemoveAll(w.Folder()); err != nil {
			panic(err)
		}
		fmt.Printf("%s is removed\n", w.Folder())
	}

	// {
	// 	opts := []fw.Option{
	// 		fw.OptID(""),
	// 		fw.OptFormat("json"),
	// 		fw.OptName(""),
	// 		fw.OptWatcher("", "json", "100ms", false, false, ""),
	// 	}
	// 	freader, err := fw.NewFileWatcher(opts...)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	freader.StartWait(prepare, cleanup)
	// }

	{
		optsFW := []fw.Option{
			fw.OptID(""),
			fw.OptFormat("json"),
			fw.OptName(""),
			fw.OptWatcher("", "json", "100ms", false, false, ""),
		}
		freader, err := fw.NewFileWatcher(optsFW...)
		Check(err)

		opts := []Option{
			OptNatsHostName(""),
			OptNatsPort(0),
			OptNatsStream("STREAM-1"),
			OptNatsStreamSubjects("STREAM-1.*"),
			OptSubject("STREAM-1.test"),
			OptKeyValue("Provider", "test-provider"),
			OptKeyValue("provider-1", "test-provider-1"), // test, should not be meta out
		}
		n3r, err := NewNats4Reader(opts...)
		Check(err)

		freader.Event = NewN3ReaderEvent(n3r)
		freader.StartWait(prepare, cleanup)
	}
}
