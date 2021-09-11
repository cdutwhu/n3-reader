package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	. "github.com/cdutwhu/n3-reader"
	fw "github.com/cdutwhu/n3-reader/file-watcher"
	cp "github.com/digisan/cli-prompt"
	"github.com/pkg/errors"
)

var mc map[string]interface{}
var err error

// use outter mc
func S(name string) string {
	return mc[name].(string)
}
func B(name string) bool {
	return mc[name].(bool)
}
func I(name string) int {
	return int(mc[name].(float64))
}

func main() {

	configPtr := flag.String("c", "./config.json", "config(json) file path")
	flag.Parse()

	mc, err = cp.PromptConfig(*configPtr)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "Invalid config file as JSON format"))
	}

	if mc != nil {
		fmt.Println("Running...")
	}

	// ------------------------------------------ //

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
			fw.OptID(S("ID")),
			fw.OptFormat(S("Format")),
			fw.OptName(S("ReaderName")),
			fw.OptWatcher(S("WatchFolder"), "", S("Interval"), B("Recursive"), B("InclHidden"), S("Ignore")),
		}
		freader, err := fw.NewFileWatcher(optsFW...)
		Check(err)

		opts := []Option{
			OptNatsHost(S("NatsHost")),
			OptNatsPort(I("NatsPort")),
			OptStream(S("Stream")),
			OptStreamSubjects(S("Stream") + ".*"),
			OptSubject(S("Stream") + "." + S("Subject")),

			// extention use for otf-reader, etc...
			OptKeyValue("Provider", "test-provider"),
			OptKeyValue("provider-1", "test-provider-1"), // test, should not be meta out
		}
		n3r, err := NewNats4Reader(opts...)
		Check(err)

		freader.Event = NewN3ReaderEvent(n3r)
		freader.StartWait(prepare, cleanup)
	}
}
