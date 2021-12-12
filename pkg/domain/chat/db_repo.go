package chat

import m "github.com/mj23978/chat-backend/pkg/models"

type ChatDBRepo interface {
	Create(chat *m.Chat) error
	Delete(filter map[string]interface{}) error
	Update(filter map[string]interface{}, fields map[string]interface{}) error
	FindBy(filter map[string]interface{}) (*m.Chat, error)
}
