package chat

import (
	m "github.com/mj23978/chat-backend/pkg/models"
)

type ChatCacheRepo interface {
	Create(chat *m.Chat) error
	Delete(id string, fields ...string) error
	Update(id string, fields map[string]interface{}) error
	Get(id string, fields ...string) (*m.Chat, error)
	GetAll(id string) (*m.Chat, error)
	Incr(id string, field uint8) (int64, error)
	Publish(id string, message interface{}) error
}
