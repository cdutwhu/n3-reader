package n3reader

import (
	"fmt"
	"os"

	stan "github.com/nats-io/stan.go"
)

type Nats4Reader struct {
	host        string
	port        int
	cluster     string
	topic       string
	sc          stan.Conn
	nConcurrent int
}

func (n4r *Nats4Reader) meta() string {
	return fmt.Sprintf(`{
		"NatsHost": "%s",
		"NatsPort": "%5d",
		"NatsCluster": "%s",
		"PublishTopic":"%s",
		"ConcurrentFiles":"%5d"
	}`, n4r.host, n4r.port, n4r.cluster, n4r.topic, n4r.nConcurrent)
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

	_, err := os.Open(fileName)
	if err != nil {
		return err
	}

	return nil
}
