package user

import m "github.com/mj23978/chat-backend/pkg/models"

type UserSerializer interface {
	Decode(input []byte) (*m.User, error)
	Encode(input *m.User) ([]byte, error)
	DecodeMap(input []byte) (map[string]interface{}, error)
	EncodeMap(input map[string]interface{}) ([]byte, error)
}