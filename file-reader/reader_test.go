package filereader

import (
	"fmt"
	"os"
	"testing"
)

func TestNewFileReader(t *testing.T) {
	opts := []Option{
		OptID(""),
		OptFormat("json"),
		OptName("Reader"),
		OptWatcher("", "json", "1s", false, false, ""),
	}
	fr, err := NewFileReader(opts...)
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
