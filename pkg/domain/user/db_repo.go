package user

import m "github.com/mj23978/chat-backend/pkg/models"

type UserDBRepo interface {
  Create(user *m.User) error
  Delete(filter map[string]interface{}) error
  Update(filter map[string]interface{}, fields map[string]interface{}) error
  FindBy(filter map[string]interface{}) (*m.User, error)
  EncryptPassword(password string) (string, error)
  DecryptPassword(hashedPass, password string) error
}
