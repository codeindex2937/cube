package message

import (
	"bytes"
	"encoding/json"
	"net/http"

	"cube/config"
	"cube/lib/database"
	"cube/lib/logger"

	"gorm.io/gorm"
)

type IService interface {
	Send(db *gorm.DB, userID int, regID uint64, ctx map[string]interface{})
}

var impl IService
var log = logger.Log

type serviceImpl struct{}

func Service() IService {
	if impl == nil {
		impl = &serviceImpl{}
	}
	return impl
}

func SetService(newImpl IService) IService {
	originImpl := impl
	impl = newImpl
	return originImpl
}

func (s *serviceImpl) Send(db *gorm.DB, userID int, regID uint64, ctx map[string]interface{}) {
	copiedCtx := make(map[string]interface{})
	for k, v := range ctx {
		copiedCtx[k] = v
	}

	if regID < 1 {
		copiedCtx["user_ids"] = []int{userID}
		requestByte, _ := json.Marshal(copiedCtx)

		resp, _ := http.Post(
			"https://chat.synology.com/webapi/entry.cgi?api=SYNO.Chat.External&method=chatbot&version=2&token=%22"+config.Conf.Token+"%22",
			"application/x-www-form-urlencoded",
			bytes.NewReader(append([]byte("payload="), requestByte...)),
		)

		if resp.StatusCode != 200 {
			log.Errorf("send message: %v", resp)
		}
	} else {
		var reg database.Registration

		tx := db.First(&reg, regID)
		if tx.Error == gorm.ErrRecordNotFound {
			log.Errorf("no record for %v", userID)
		}

		requestByte, _ := json.Marshal(copiedCtx)

		resp, _ := http.Post(
			"https://chat.synology.com/webapi/entry.cgi?api=SYNO.Chat.External&method=incoming&version=2&token="+reg.Token,
			"application/x-www-form-urlencoded",
			bytes.NewReader(append([]byte("payload="), requestByte...)),
		)

		if resp.StatusCode != 200 {
			log.Errorf("send message: %v", resp)
		}
	}
}
