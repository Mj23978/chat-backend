package model

type Chat struct {
	ID            string    `json:"channel"  msgpack:"channel"  validate:"empty=false"`
	Name          string    `json:"name"  msgpack:"name"  validate:"empty=false"`
	Members       []string  `json:"members"  msgpack:"members"`
	TypingMembers []string  `json:"typing_members"  msgpack:"typing_members"`
	Messages      []Message `json:"messages"  msgpack:"messages"`
	ImageUrl      string    `json:"image_url"  msgpack:"image_url"`
}

type Command struct {
	Name        string   `json:"namr"  msgpack:"namr"  validate:"empty=false"`
	Description string   `json:"description"  msgpack:"description" `
	Args        []string `json:"args"  msgpack:"args"`
}

type Group struct {
	ID              string   `json:"channel"  msgpack:"channel"  validate:"empty=false"`
	GroupName       string   `json:"group_name"  msgpack:"group_name"  validate:"empty=false"`
	GroupImage      string   `json:"group_image"  msgpack:"group_image"`
	LastMessageTime string   `json:"last_message_time"  msgpack:"last_message_time"`
	MembersIds      []string `json:"members_ids"  msgpack:"members_ids"`
}
