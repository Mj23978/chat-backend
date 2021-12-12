package model

import "time"

type Health struct {
	Username string    `json:"username"`
	Time     time.Time `json:"time"`
}

type GameServer struct {
	Type     string      `json:"type"`
	HasError bool        `json:"has_error"`
	Payload  interface{} `json:"payload"`
}