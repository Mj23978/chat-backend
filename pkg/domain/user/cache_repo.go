package user

import m "github.com/mj23978/chat-backend/pkg/models"

type UserCacheRepo interface {
	Create(user *m.User) error
	Delete(id string, fields ...string) error
	Update(id string, fields map[string]interface{}) error
	Get(id string, fields ...string) (*m.User, error)
}
