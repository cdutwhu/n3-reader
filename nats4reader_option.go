package n3reader

type Option func(*Nats4Reader) error

func (n3r *Nats4Reader) setOption(options ...Option) error {
	for _, opt := range options {
		if err := opt(n3r); err != nil {
			return err
		}
	}
	return nil
}
