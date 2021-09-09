package n3reader

import (
	"encoding/json"
	"fmt"
	"log"
	"unicode"

	jt "github.com/digisan/json-tool"
	"github.com/nats-io/nats.go"
)

type Nats4Reader struct {
	host           string                // option,    meta
	port           int                   // option,    meta
	stream         string                // option,    meta
	streamSubjects string                // option,    no meta
	subject        string                // option,    meta
	nc             *nats.Conn            // no option, no meta
	js             nats.JetStreamContext // no option, no meta

	// for outter user filling, only first UpperCase key can be meta
	kvInfo map[string]interface{} // option as kv, Upper-Key meta
}

func (n4r *Nats4Reader) meta() string {

	// keep an eye on last comma
	return fmt.Sprintf(`{
		"NatsHost": "%s",
		"NatsPort": "%5d",
		"Stream": "%s",
		"Subject": "%s"
	}`,
		n4r.host,
		n4r.port,
		n4r.stream,
		n4r.subject,
	)
}

func (n4r *Nats4Reader) metaWithKV() string {

	// only select upper case key kv for meta string
	m := make(map[string]interface{})
	for k, v := range n4r.kvInfo {
		if unicode.IsUpper(rune(k[0])) {
			m[k] = v
		}
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		log.Fatalln(err)
	}
	return jt.MergeSgl(n4r.meta(), string(bytes))
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
	n4r := &Nats4Reader{kvInfo: make(map[string]interface{})}
	if err := n4r.setOption(options...); err != nil {
		return nil, err
	}
	return n4r, n4r.initNatsJS()
}
