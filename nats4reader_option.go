package n3reader

import (
	"github.com/pkg/errors"
)

const (
	dfltHost       = "localhost"    // nats default
	dfltPort       = 4222           // nats default
	dfltCluster    = "test-cluster" // nats default
	dfltConcurrent = 10             // safe default
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

//
// set the nats server name or ip address
// empty string will result in localhost as defalt hostname
//
func NatsHostName(hostName string) Option {
	return func(n4r *Nats4Reader) error {
		if hostName != "" {
			n4r.host = hostName
		}
		n4r.host = dfltHost // nats default
		return nil
	}
}

//
// set the nats server communication port
// port value of 0 or less will result in default nats port 4222
//
func NatsPort(port int) Option {
	return func(n4r *Nats4Reader) error {
		if port > 0 {
			n4r.port = port
			return nil
		}
		n4r.port = dfltPort
		return nil
	}
}

//
// set the nats streaming server cluster name.
// empty string will result in nats default of 'test-cluster'
//
func NatsClusterName(clusterName string) Option {
	return func(n4r *Nats4Reader) error {
		if clusterName != "" {
			n4r.cluster = clusterName
		}
		n4r.cluster = dfltCluster
		return nil
	}
}

//
// set the name of the nats topic to publish data once parsed
// from the input files
//
func Topic(tName string) Option {
	return func(n4r *Nats4Reader) error {
		if tName == "" {
			return errors.New("must have Topic (nats topic to which reader will publish parsed data).")
		}

		// topic regex check
		ok, err := ValidateNatsTopic(tName)
		if ok {
			n4r.topic = tName
			return nil
		}
		return errors.Wrap(err, "Topic option error")
	}
}

//
// set the number of input files that can be handled concurrently
// set if number of filehandles on OS is a problem
// defaults to 10
//
func ConcurrentFiles(n int) Option {
	return func(n4r *Nats4Reader) error {
		if n == 0 {
			n4r.nConcurrent = dfltConcurrent // safe default
			return nil
		}
		n4r.nConcurrent = n
		return nil
	}
}
