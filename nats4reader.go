package n3reader

import (
	"encoding/json"
	"fmt"
	"unicode"

	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/nats-io/nats.go"
)

type NatsReader struct {
	host           string                // option,    no meta
	port           int                   // option,    no meta
	stream         string                // option,    meta
	streamSubjects string                // option,    no meta
	subject        string                // option,    meta
	nc             *nats.Conn            // no option, no meta
	js             nats.JetStreamContext // no option, no meta

	// for outter user filling, only first UpperCase key can be meta
	kvInfo map[string]interface{} // option as kv, Upper-Key meta
}

func (nr *NatsReader) meta() string {
	m := map[string]interface{}{
		"Stream":  nr.stream,
		"Subject": nr.subject,
	}
	data, err := json.Marshal(m)
	lk.FailOnErr("%v", err)
	return string(data)
}

func (nr *NatsReader) exMeta() string {

	// only select upper case key kv for meta string
	m := make(map[string]interface{})
	for k, v := range nr.kvInfo {
		if unicode.IsUpper(rune(k[0])) {
			m[k] = v
		}
	}
	data, err := json.Marshal(m)
	lk.FailOnErr("%v", err)
	return jt.MergeSgl(nr.meta(), string(data))
}

func (nr *NatsReader) initNatsJS() (err error) {

	// create connection & JetStreamContext
	nr.nc, err = nats.Connect(fmt.Sprintf("nats://%s:%d", nr.host, nr.port)) // "nats://127.0.0.1:4222"
	if err != nil {
		return err
	}
	nr.js, err = nr.nc.JetStream()
	if err != nil {
		return err
	}

	// check if the stream already exists; if not, create it
	stream, err := nr.js.StreamInfo(nr.stream)

	lk.LogWhen(err != nil, "%v", err) // notice 'not found, ready to create a new one'

	if stream == nil {
		lk.Log("creating stream %q and subjects %q", nr.stream, nr.streamSubjects)
		_, err = nr.js.AddStream(
			&nats.StreamConfig{
				Name:     nr.stream,
				Subjects: []string{nr.streamSubjects},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewNats4Reader(options ...Option) (*NatsReader, error) {
	nr := &NatsReader{kvInfo: make(map[string]interface{})}
	if err := nr.setOption(options...); err != nil {
		return nil, err
	}
	return nr, nr.initNatsJS()
}
