package message

import (
	"github.com/go-redis/redis"
	m "github.com/mj23978/chat-backend/pkg/models"
)

type MessageCacheRepo interface {
	Create(message *m.Message) error
	Delete(id string, fields ...string) error
	Update(id string, fields map[string]interface{}) error
	Get(id string, fields ...string) (*m.Message, error)
	GetAll(id string) (*m.Message, error)
	PubSub(channels ...string) (*redis.PubSub, error)
}
