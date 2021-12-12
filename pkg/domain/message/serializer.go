package message

import m "github.com/mj23978/chat-backend/pkg/models"

type MessageSerializer interface {
	Decode(input []byte) (*m.Message, error)
	Encode(input *m.Message) ([]byte, error)
}
