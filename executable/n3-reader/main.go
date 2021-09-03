package main

import (
	"fmt"
	"os"

	fr "github.com/cdutwhu/n3-reader/file-reader"
)

func main() {
	opts := []fr.Option{
		fr.OptID(""),
		fr.OptFormat("json"),
		fr.OptName(""),
		fr.OptWatcher("", "json", "100ms", false, false, ""),
	}
	fr, err := fr.NewFileReader(opts...)
	if err != nil {
		panic(err)
	}
	cleanup := func(folder string) {
		if err := os.RemoveAll(folder); err != nil {
			panic(err)
		}
		fmt.Printf("%s is removed\n", folder)
	}
	fr.StartWait(cleanup)
}
