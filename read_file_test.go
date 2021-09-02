package n3reader

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestNewFileReader(t *testing.T) {

	opts := []Option{
		Watcher("", "json", "1s", false, false, ""),
	}
	fr, err := NewFileReader(opts...)
	if err == nil {

		// signal handler for shutdown
		closed := make(chan struct{})
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			fmt.Println("\nreader shutting down")
			fr.Close()
			fmt.Println("otf-reader closed")
			close(closed)
		}()

		fr.Start()
		<-closed
	}
}
