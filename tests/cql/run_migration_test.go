package migration

import (
	m "github.com/mj23978/chat-backend/pkg/models"
	"github.com/mj23978/chat-backend/pkg/repository/cassandra"
	"testing"
)

func TestExample(t *testing.T) {
	const ks = "chat"

	// Create keyspace
	cas := cassandra.NewCassandra(cassandra.Config{Address: "127.0.0.1", Keyspace: ks, Username: "", Password: ""})

	err := cas.CreateKeyspace(ks)
	//err = cas.CreateType("StatusType", m.StatusType(0))
	//if err != nil {
	//	panic(err)
	//}
	//err = cas.CreateType("User", m.User{})
	//if err != nil {
	//	panic(err)
	//}
	err = cas.CreateType("Attachment", m.Attachment{})
	if err != nil {
		panic(err)
	}
	err = cas.CreateType("Message", m.Message{})
	if err != nil {
		panic(err)
	}
	userTable := cassandra.Table{Indexes: []string{"username", "email"}, Name: "users", Type: m.User{}, PrimaryKeys: [][]string{[]string{"id"}, []string{"email", "username"}}}
	err = userTable.CreateTable(cas)
	if err != nil {
		panic(err)
	}
	chatTable := cassandra.Table{Name: "chats", Type: m.Chat{}, PrimaryKeys: [][]string{[]string{"id"}}}
	err = chatTable.CreateTable(cas)
	if err != nil {
		panic(err)
	}
}
