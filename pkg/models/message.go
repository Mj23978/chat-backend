package model

type Message struct {
	ID             string              `json:"id"  msgpack:"id"  validate:"empty=false"`
	Text           string              `json:"text"  msgpack:"text"  validate:"empty=false"`
	Status         MessageStatus       `json:"status"  msgpack:"status"  validate:"empty=false"`
	From           string              `json:"from"  msgpack:"from"  validate:"empty=false"`
	CreatedAt      string              `json:"created_at"  msgpack:"created_at"  validate:"empty=false"`
	UpdatedAt      string              `json:"updated_at"  msgpack:"updated_at"  validate:"empty=false"`
	User           User                `json:"user"  msgpack:"user"`
	Silent         bool                `json:"silent"  msgpack:"silent"`
	Pinned         bool                `json:"pinned"  msgpack:"pinned"`
	Likes          int32               `json:"likes"  msgpack:"likes"`
	Views          int32               `json:"views"  msgpack:"views"`
	Tags           []string            `json:"tags"  msgpack:"tags"`
	Content        []string            `json:"content"  msgpack:"content"`
	Attachments    []Attachment        `json:"attachments"  msgpack:"attachments"`
	MentionedUsers []User              `json:"mentioned_users"  msgpack:"mentioned_users"`
	ReplyTo        string              `json:"reply_to"  msgpack:"reply_to"`
	ReplyCount     int32               `json:"reply_count"  msgpack:"reply_count"`
	To             string              `json:"to"  msgpack:"to"`
	Command        string              `json:"command"  msgpack:"command"`
	ParentID       string              `json:"parent_id"  msgpack:"parent_id"`
	ExtraData      map[string]interface{} `json:"extra_data"  msgpack:"extra_data"`
}

type Reaction struct {
	ID        string              `json:"id"  msgpack:"id"  validate:"empty=false"`
	User      User                `json:"user"  msgpack:"user"`
	MessageId string              `json:"message_id"  msgpack:"message_id"  `
	UserId    string              `json:"user_id"  msgpack:"user_id"  `
	CreatedAt string              `json:"created_at"  msgpack:"created_at"`
	ExtraData map[string]interface{} `json:"extra_data"  msgpack:"extra_data"`
}

type Story struct {
	ID          string   `json:"id"  msgpack:"id"  validate:"empty=false"`
	PublishedAt string   `json:"created_at"  msgpack:"created_at"  validate:"empty=false"`
	User        User     `json:"user"  msgpack:"user"  validate:"empty=false"`
	IsSeen      bool     `json:"is_seen"  msgpack:"is_seen"  validate:"empty=false"`
	Images      []string `json:"images"  msgpack:"images"  validate:"empty=false"`
}
