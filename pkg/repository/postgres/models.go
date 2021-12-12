package postgres

import (
	"database/sql"
	m "github.com/mj23978/chat-backend/pkg/models"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          string         `json:"id"  msgpack:"id" gorm:"primaryKey" mapstructure:"id"`
	Age         uint8          `json:"age" msgpack:"age" mapstructure:"age"`
	ActivatedAt sql.NullTime   `json:"activated_at" msgpack:"activated_at" mapstructure:"activated_at"`
	Phone       string         `json:"phone" msgpack:"phone" mapstructure:"phone"`
	DeletedAt   gorm.DeletedAt `gorm:"index" mapstructure:"deleted_at"`
	Username    string         `json:"username" msgpack:"username" validate:"empty=false" gorm:"index:,unique" mapstructure:"username"`
	DisplayName string         `json:"display_name" msgpack:"display_name" validate:"empty=false" mapstructure:"display_name"`
	Password    string         `json:"password" msgpack:"password" mapstructure:"password"`
	Token       string         `json:"token" msgpack:"token" mapstructure:"token"`
	PhotoUrl    string         `json:"photo_url" msgpack:"photo_url" mapstructure:"photo_url"`
	Email       string         `json:"email" msgpack:"email" gorm:"index:,unique;not null" mapstructure:"email"`
	Role        string         `json:"role" msgpack:"role" mapstructure:"role"`
	LastActive  time.Time      `json:"last_active" msgpack:"last_active" mapstructure:"last_active"`
	Status      m.StatusType   `json:"status" msgpack:"status" gorm:"not null" mapstructure:"status"`
	Friends     []*User        `json:"friends" msgpack:"friends" gorm:"many2many:user_friends;foreignKey:Username;joinForeignKey:UsernameID;References:Username;joinReferences:FriendID" mapstructure:"friends"`
	CreatedAt   time.Time      `json:"created_at" msgpack:"created_at" mapstructure:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" msgpack:"updated_at" mapstructure:"updated_at"`
}
