package model

type StatusType int32

const (
	Online  StatusType = iota
	Away    StatusType = iota
	Playing StatusType = iota
	Offline StatusType = iota
)

type User struct {
	ID          string     `json:"id"  msgpack:"id"`
	Username    string     `json:"username" msgpack:"username" validate:"empty=false"`
	DisplayName string     `json:"display_name" msgpack:"display_name" validate:"empty=false"`
	Password    string     `json:"password" msgpack:"password"`
	Token       string     `json:"token" msgpack:"token"`
	PhotoUrl    string     `json:"photo_url" msgpack:"photo_url"`
	Email       string     `json:"email" msgpack:"email"`
	Role        string     `json:"role" msgpack:"role"`
	LastActive  string     `json:"last_active" msgpack:"last_active"`
	Status      StatusType `json:"status" msgpack:"status"`
	Friends     []string   `json:"friends" msgpack:"friends"`
	CreatedAt   string     `json:"created_at" msgpack:"created_at"`
	UpdatedAt   string     `json:"updated_at" msgpack:"updated_at"`
}

type UserPublic struct {
	Username string     `json:"username" msgpack:"username" validate:"empty=false"`
	PhotoUrl string     `json:"photo_url" msgpack:"photo_url"`
	Status   StatusType `json:"status" msgpack:"status"`
	//Status   int32  `json:"status" msgpack:"status"`
}

type Member struct {
	Invited          bool   `json:"invited" msgpack:"invited"`
	Banned           bool   `json:"banned" msgpack:"banned"`
	ShadowBanned     bool   `json:"shadow_banned" msgpack:"shadow_banned"`
	IsModerator      bool   `json:"is_moderator" msgpack:"is_moderator"`
	User             User   `json:"user" msgpack:"user"`
	InviteAcceptedAt string `json:"invite_accepted_at" msgpack:"invite_accepted_at"`
	InviteRejectedAt string `json:"invite_rejected_at" msgpack:"invite_rejected_at"`
	Role             string `json:"role" msgpack:"role"`
	CreatedAt        string `json:"created_at" msgpack:"created_at"`
	UpdatedAt        string `json:"updated_at" msgpack:"updated_at"`
}
