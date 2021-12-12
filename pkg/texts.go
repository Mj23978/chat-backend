package test


const (
	// Broker Topics
	BrokerGetUser = "get-user"
	BrokerGameChannel = "game-channel"
	
	// Emitter Topics
	EmitterRequestGetUser = "req:get-user"
)

type Msg struct {
	test string
}