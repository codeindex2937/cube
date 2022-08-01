package event

import (
	"cube/lib/logger"
)

var log = logger.Log()

type Subscriber func(event string, ctx interface{})

type IService interface {
	Subscribe(e string, f Subscriber)
	Publish(e string, ctx interface{})
}

type ServiceImpl struct {
	subscriberMap map[string][]Subscriber
}

func NewService() IService {
	return &ServiceImpl{
		subscriberMap: make(map[string][]Subscriber),
	}
}

func (s *ServiceImpl) Subscribe(e string, f Subscriber) {
	subscribers, ok := s.subscriberMap[e]
	if !ok {
		subscribers = make([]Subscriber, 0)
		s.subscriberMap[e] = subscribers
	}
	s.subscriberMap[e] = append(subscribers, f)
}

func (s *ServiceImpl) Publish(e string, ctx interface{}) {
	subscribers, ok := s.subscriberMap[e]
	if !ok {
		log.Warnf("unknown event: %s", e)
		return
	}

	for _, s := range subscribers {
		s(e, ctx)
	}
}
