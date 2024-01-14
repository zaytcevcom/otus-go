package sender

type Sender struct {
	logger Logger
	broker MessageBroker
	doneCh chan interface{}
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type MessageBroker interface {
	Subscribe(handler func(body []byte) error) error
}

func New(logger Logger, broker MessageBroker) *Sender {
	return &Sender{
		logger: logger,
		broker: broker,
		doneCh: make(chan interface{}),
	}
}

func (s Sender) Start() error {

	s.logger.Debug("Sender started!")

	err := s.broker.Subscribe(func(body []byte) error {
		s.logger.Info(string(body))
		return nil
	})
	if err != nil {
		return err
	}

	for range s.doneCh {
		return nil
	}

	return nil
}

func (s Sender) Stop() {
	close(s.doneCh)
}
