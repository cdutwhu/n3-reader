package n3reader

import (
	"testing"

	fr "github.com/cdutwhu/n3-reader/file-reader"
)

func TestNewN3Reader(t *testing.T) {
	if n3r, err := NewNats4Reader(); err == nil {
		opts := []fr.Option{
			fr.OptID(""),
			fr.OptFormat("json"),
			fr.OptName(""),
			fr.OptWatcher("", "json", "100ms", false, false, ""),
		}
		if freader, err := fr.NewFileReader(opts...); err == nil {
			freader.Event = NewN3ReaderEvent(n3r)
			freader.StartWait()
		}
	}
}
