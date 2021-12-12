package model

type Channel struct {
	ID string `json:"id"  msgpack:"id"  validate:"empty=false"`
}

type Room struct {
	ID string `json:"id"  msgpack:"id"  validate:"empty=false"`
}

type Tag struct {
	ID   string `json:"id"  msgpack:"id"  validate:"empty=false"`
	Text string `json:"text"  msgpack:"text"  validate:"empty=false"`
}
