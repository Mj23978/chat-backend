package model

type ChannelMessage struct {
	Type     string      `json:"type"`
	From     string      `json:"from"`
	HasError bool        `json:"has_error"`
	Payload  interface{} `json:"payload"`
}

type ErrorPayload struct {
	Error    string      `json:"error"`
	ID       string      `json:"id"`
	Payload  interface{} `json:"payload"`
}

type LocalPayload struct {
	HasError bool        `json:"has_error"`
	ID       string      `json:"id"`
	Payload  interface{} `json:"payload"`
}

const (
	GameJoin        = "join"
	GameLeave       = "leave"
	GamePublish     = "publish"
	GameSubscribe   = "subscribe"
	GameUnSubscribe = "unsubscribe"
	GameBroadcast   = "broadcast"
	GameHealth      = "health"
)