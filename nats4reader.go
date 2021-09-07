package n3reader

import (
	"context"
	"fmt"
	"log"
	"os"

	jt "github.com/digisan/json-tool"
	"github.com/nats-io/nats.go"
)

type Nats4Reader struct {
	host           string                // option, meta
	port           int                   // option, meta
	stream         string                // option, meta
	streamSubjects string                // option, no meta
	subject        string                // option meta
	nc             *nats.Conn            // no option, no meta
	js             nats.JetStreamContext // no option, no meta
}

func (n4r *Nats4Reader) meta() string {
	return fmt.Sprintf(`{
		"NatsHost": "%s",
		"NatsPort": "%5d",
		"StreamName": "%s",
		"Subject": "%s",
	}`,
		n4r.host,
		n4r.port,
		n4r.stream,
		n4r.subject,
	)
}

func (n4r *Nats4Reader) initNatsJS() (err error) {

	// create connection & JetStreamContext
	n4r.nc, err = nats.Connect(fmt.Sprintf("nats://%s:%d", n4r.host, n4r.port)) // "nats://127.0.0.1:4222"
	if err != nil {
		return err
	}
	n4r.js, err = n4r.nc.JetStream()
	if err != nil {
		return err
	}

	// check if the stream already exists; if not, create it
	stream, err := n4r.js.StreamInfo(n4r.stream)
	if err != nil {
		log.Println(err) // notice 'not found, ready to create a new one'
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", n4r.stream, n4r.streamSubjects)
		_, err = n4r.js.AddStream(
			&nats.StreamConfig{
				Name:     n4r.stream,
				Subjects: []string{n4r.streamSubjects},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewNats4Reader(options ...Option) (*Nats4Reader, error) {
	n4r := &Nats4Reader{}
	if err := n4r.setOption(options...); err != nil {
		return nil, err
	}
	return n4r, n4r.initNatsJS()
}

func (n4r *Nats4Reader) PubAsJSON(fileName, meta string) error {
	fmt.Println("Publishing:", fileName)

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cOut, jarr := jt.ScanObject(ctx, f, false, true, jt.OUT_ORI)
	if !jarr {
		log.Println("not json array")
	}
	for result := range cOut {
		if result.Err != nil {
			log.Println(result.Err)
			return result.Err
		}

		data := jt.Minimize(fmt.Sprintf(`{"meta":%s, "data":%s}`, meta, result.Obj), true)
		ack, err := n4r.js.Publish(n4r.subject, []byte(data))
		if err != nil {
			return err
		}
		fmt.Println("ACK:", ack.Stream)
	}

	return nil
}
