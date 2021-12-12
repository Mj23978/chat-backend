package gameNode

import (
	"fmt"

	nprotoo "github.com/mj23978/chat-backend-x/broker/nats"
	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	test "github.com/mj23978/chat-backend/pkg"
)

func (gn *gameNode) handleRequest() {
	gn.protoo.OnRequest(gn.ServiceNode.GetRPCChannel(), func(request nprotoo.Request, accept nprotoo.RespondFunc, reject nprotoo.RejectFunc) {
		method := request.Method
		data := request.Data
		log.Debugf("handleRequest: method => %s, data => %v", method, data)

		var result interface{}
		err := NewNpError(400, fmt.Sprintf("Unkown method [%s]", method))

		switch method {
		case test.BrokerGetUser:
			var msgData test.Msg
			if err = data.Unmarshal(&msgData); err == nil {
				result, err = publish(msgData)
			}
		case test.BrokerGameChannel:
			var msgData test.Msg
			if err = data.Unmarshal(&msgData); err == nil {
				result, err = publish(msgData)
			}
		}

		if err != nil {
			reject(err.Code, err.Reason)
		} else {
			accept(result)
		}
	})
}

func NewNpError(code int, reason string) *nprotoo.Error {
	err := nprotoo.Error{
		Code:   code,
		Reason: reason,
	}
	return &err
}

func publish(msg test.Msg) (string, *nprotoo.Error) {
	return fmt.Sprintf("Good - %v", msg), nil
}
