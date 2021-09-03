package filereader

import (
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

	fr.StartWait()
}
