package n3reader

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	dfltHost           = "127.0.0.1"
	dfltPort           = 4222
	dfltStream         = "STREAM-1"
	dfltStreamSubjects = "STREAM-1.*"
	dfltSubject        = "STREAM-1.sub1"
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

func OptNatsHostName(hostName string) Option {
	return func(n4r *Nats4Reader) error {
		return SetIfNotEmpty(&n4r.host, hostName, dfltHost)
	}
}

func OptNatsPort(port int) Option {
	return func(n4r *Nats4Reader) error {
		return SetIfNotZero(&n4r.port, port, dfltPort)
	}
}

func OptNatsStream(stream string) Option {
	return func(n4r *Nats4Reader) error {
		return SetIfNotEmpty(&n4r.stream, stream, dfltStream)
	}
}

func OptNatsStreamSubjects(streamSubjects string) Option {
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
