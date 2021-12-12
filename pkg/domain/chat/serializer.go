package chat

import m "github.com/mj23978/chat-backend/pkg/models"

type ChatSerializer interface {
	Decode(input []byte) (*m.Chat, error)
	Encode(input *m.Chat) ([]byte, error)
	DecodeMap(input []byte) (map[string]interface{}, error)
	EncodeMap(input map[string]interface{}) ([]byte, error)
}
