package n3reader

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	dfltHost           = "127.0.0.1"
	dfltPort           = 4222
	dfltStream         = "STREAM"
	dfltStreamSubjects = "STREAM.*"
	dfltSubject        = "STREAM.sub"
)

type Option func(*Nats4Reader) error

func (n3r *Nats4Reader) setOption(options ...Option) error {
	for i, opt := range options {
		if err := opt(n3r); err != nil {
			return errors.Wrap(err, fmt.Sprintf("@%d", i))
		}
	}
	return nil
}

// Options

func OptNatsHost(hostName string) Option {
	return func(n4r *Nats4Reader) error {
		return SetIfNotEmpty(&n4r.host, hostName, dfltHost)
	}
}

func OptNatsPort(port int) Option {
	return func(n4r *Nats4Reader) error {
		return SetIfNotZero(&n4r.port, port, dfltPort)
	}
}

func OptStream(stream string) Option {
	return func(n4r *Nats4Reader) error {
		return SetIfNotEmpty(&n4r.stream, stream, dfltStream)
	}
}

func OptStreamSubjects(streamSubjects string) Option {
	return func(n4r *Nats4Reader) error {
		return SetIfNotEmpty(&n4r.streamSubjects, streamSubjects, dfltStreamSubjects)
	}
}

func OptSubject(subject string) Option {
	return func(n4r *Nats4Reader) error {
		validate := func(s string) (bool, error) {
			if s == "" {
				return false, errors.New("must have Subject (nats subject to which reader will publish parsed data)")
			}
			return ValidateNatsSubject(subject)
		}
		return SetIfValidStr(&n4r.subject, subject, validate)
	}
}

//////////////////////////////////////////////////

// for outter user like otf-reader
func OptKeyValue(key string, value interface{}) Option {
	return func(n4r *Nats4Reader) error {
		n4r.kvInfo[key] = value
		return nil
	}
}
