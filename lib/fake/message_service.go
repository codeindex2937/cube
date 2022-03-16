package fake

import (
	"cube/lib/logger"
)

var log = logger.Log

type MessageService struct{}

func (s *MessageService) Send(userID, message string) {
	log.Infof("%v %v\n", userID, message)
}
