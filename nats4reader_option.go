package n3reader

import (
	"fmt"

	"github.com/pkg/errors"
)

// If Invalid value set in config, below values apply
const (
	Host           = "127.0.0.1"
	Port           = 4222
	Stream         = "STREAM"
	StreamSubjects = "STREAM.*"
	Subject        = "STREAM.sub"
)

type Option func(*NatsReader) error

func (nr *NatsReader) setOption(options ...Option) error {
	for i, opt := range options {
		if err := opt(nr); err != nil {
			return errors.Wrap(err, fmt.Sprintf("@%d", i))
		}
	}
	return nil
}

// Options

func OptNatsHost(hostName string) Option {
	return func(nr *NatsReader) error {
		return SetIfNotEmpty(&nr.host, hostName, Host)
	}
}

func OptNatsPort(port int) Option {
	return func(nr *NatsReader) error {
		return SetIfNotZero(&nr.port, port, Port)
	}
}

func OptStream(stream string) Option {
	return func(nr *NatsReader) error {
		return SetIfNotEmpty(&nr.stream, stream, Stream)
	}
}

func OptStreamSubjects(streamSubjects string) Option {
	return func(nr *NatsReader) error {
		return SetIfNotEmpty(&nr.streamSubjects, streamSubjects, StreamSubjects)
	}
}

func OptSubject(subject string) Option {
	return func(nr *NatsReader) error {
		validate := func(s string) (bool, error) {
			if s == "" {
				return false, errors.New("must have Subject (nats subject to which reader will publish parsed data)")
			}
			return ValidateNatsSubject(subject)
		}
		return SetIfValidStr(&nr.subject, subject, validate)
	}
}

//////////////////////////////////////////////////

// for outter user like otf-reader
func OptKeyValue(key string, value interface{}) Option {
	return func(nr *NatsReader) error {
		nr.kvInfo[key] = value
		return nil
	}
}
