package n3reader

import (
	"github.com/pkg/errors"
)

const (
	dfltHost           = "127.0.0.1"           // nats default
	dfltPort           = 4222                  // nats default
	dfltStream         = "TEST-STREAM"         // nats default
	dfltStreamSubjects = "TEST-STREAM.*"       // nats default
	dfltSubject        = "TEST-STREAM.created" // nats default
	dfltConcurrent     = 10                    // safe default
)

type Option func(*Nats4Reader) error

func (n3r *Nats4Reader) setOption(options ...Option) error {
	for _, opt := range options {
		if err := opt(n3r); err != nil {
			return err
		}
	}
	return nil
}

func OptNatsHostName(hostName string) Option {
	return func(n4r *Nats4Reader) error {
		SetIfNotEmpty(&n4r.host, hostName, dfltHost)
		return nil
	}
}

func OptNatsPort(port int) Option {
	return func(n4r *Nats4Reader) error {
		SetIfNotZero(&n4r.port, port, dfltPort)
		return nil
	}
}

func OptNatsStream(stream string) Option {
	return func(n4r *Nats4Reader) error {
		SetIfNotEmpty(&n4r.stream, stream, dfltStream)
		return nil
	}
}

func OptNatsStreamSubjects(streamSubjects string) Option {
	return func(n4r *Nats4Reader) error {
		SetIfNotEmpty(&n4r.streamSubjects, streamSubjects, dfltStreamSubjects)
		return nil
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

func OptConcurrentFiles(n int) Option {
	return func(n4r *Nats4Reader) error {
		SetIfNotZero(&n4r.nConcurrent, n, dfltConcurrent)
		return nil
	}
}
