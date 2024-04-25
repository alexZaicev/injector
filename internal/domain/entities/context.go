package entities

import "regexp"

// MessageBrokerContext is the actual state of the message broker, holding all the configured
// queues, topics etc.
type MessageBrokerContext struct {
	Exchanges []*Exchange
}

func NewMessageBrokerContext() *MessageBrokerContext {
	return &MessageBrokerContext{}
}

func (mbc *MessageBrokerContext) AddExchange(exchange *Exchange) {
	mbc.Exchanges = append(mbc.Exchanges, exchange)
}

func (mbc *MessageBrokerContext) FindQueueByName(name string) (*Exchange, bool) {
	for _, exchange := range mbc.Exchanges {
		if exchange.Name == name && exchange.Kind == ExchangeKindQueue {
			return exchange, true
		}
	}

	return nil, false
}

func (mbc *MessageBrokerContext) FindExchangeByName(name string) (*Exchange, bool) {
	for _, exchange := range mbc.Exchanges {
		if exchange.Name == name {
			return exchange, true
		}
	}

	return nil, false
}

func (mbc *MessageBrokerContext) FindQueuesForTopic(name string) ([]*Exchange, error) {
	topic, ok := mbc.FindExchangeByName(name)
	if !ok {
		return nil, nil
	}

	var queues []*Exchange

	for _, exchange := range mbc.Exchanges {
		if exchange.Kind != ExchangeKindQueue {
			continue
		}

		matched, err := regexp.MatchString(topic.Pattern, exchange.Name)
		if err != nil {
			return nil, err
		}

		if matched {
			queues = append(queues, exchange)
		}
	}

	return queues, nil
}
