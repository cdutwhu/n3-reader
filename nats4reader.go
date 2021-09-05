package n3reader

import (
	"fmt"
	"os"
)

type Nats4Reader struct {
	host           string
	port           int
	stream         string
	streamSubjects string
	subject        string
	nConcurrent    int
}

func (n4r *Nats4Reader) meta() string {
	return fmt.Sprintf(`{
		"NatsHost": "%s",
		"NatsPort": "%5d",
		"StreamName": "%s",
		"StreamSubjects": "%s",
		"Subject": "%s",
		"ConcurrentFiles": "%5d"
	}`,
		n4r.host,
		n4r.port,
		n4r.stream,
		n4r.streamSubjects,
		n4r.subject,
		n4r.nConcurrent,
	)
}

func NewNats4Reader(options ...Option) (*Nats4Reader, error) {
	n4r := &Nats4Reader{}
	if err := n4r.setOption(options...); err != nil {
		return nil, err
	}
	return n4r, nil
}

func (n4r *Nats4Reader) PubAsJSON(fileName, meta string) error {
	fmt.Println("Publishing:", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// jt.ScanObject()

	return nil
}
