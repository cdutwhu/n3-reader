package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	n3r "github.com/cdutwhu/n3-reader"
)

func main() {
	opts := []n3r.Option{
		n3r.OptID(""),
		n3r.OptFormat("json"),
		n3r.OptName(""),
		n3r.OptWatcher("", "json", "100ms", false, false, ""),
	}
	fr, err := n3r.NewFileReader(opts...)
	if err != nil {
		panic(err)
	}

	// signal handler for shutdown
	closed := make(chan struct{})
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nreader shutting down")
		fr.Close()
		fmt.Println("reader closed")
		close(closed)
	}()

	fr.Start()
	<-closed
}
