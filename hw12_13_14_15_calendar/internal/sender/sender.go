package sender

type Sender struct {
	logger Logger
	broker MessageBroker
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type MessageBroker interface {
	Subscribe(handler func(body []byte) error) error
	Close() error
}

func New(logger Logger, broker MessageBroker) *Sender {
	return &Sender{
		logger: logger,
		broker: broker,
	}
}

func (s Sender) Start() error {

	s.logger.Debug("Sender started!")

	return s.broker.Subscribe(func(body []byte) error {
		s.logger.Info(string(body))
		return nil
	})
}

func (s Sender) Stop() error {
	return s.broker.Close()
}
