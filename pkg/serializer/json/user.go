package json

import (
	"encoding/json"

	m "github.com/mj23978/chat-backend/pkg/models"
	"github.com/pkg/errors"
)

type User struct{}

func (u *User) Decode(input []byte) (*m.User, error) {
	user := new(m.User)
	if e := json.Unmarshal(input, user); e != nil {
		return nil, errors.Wrap(e, "serializer.Json.Decode")
	}
	return user, nil
}

func (u *User) Encode(input *m.User) ([]byte, error) {
	rawMsg, e := json.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Json.Encode")
	}
	return rawMsg, nil
}

func (u *User) DecodeMap(input []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	if e := json.Unmarshal(input, &res); e != nil {
		return res, errors.Wrap(e, "serializer.Json.DecodeMap")
	}
	return res, nil
}

func (u *User) EncodeMap(input map[string]interface{}) ([]byte, error) {
	rawMsg, e := json.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Json.EncodeMap")
	}
	return rawMsg, nil
}
