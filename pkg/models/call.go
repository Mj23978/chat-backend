package model

import "time"

type CallLog struct {
	ID          string    `json:"id"  msgpack:"id"  validate:"empty=false"`
	ReceiverPic string    `json:"receiver_pic"  msgpack:"receiver_pic" `
	CallerPic   string    `json:"caller_pic"  msgpack:"caller_pic"`
	CreatedAt   time.Time `json:"created_at"  msgpack:"created_at"  validate:"empty=false"`
}
