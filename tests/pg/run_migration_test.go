package migration

import (
	"github.com/mj23978/chat-backend/pkg/repository/postgres"
	"testing"
)

func TestExample(t *testing.T) {
	const ks = "chat"

	// Create keyspace
	pg, err := postgres.NewPostgres(&postgres.Config{Address: "127.0.0.1", Database: "yugabyte", Username: "yugabyte", Password: "", Port: "5433"})
	if err != nil {
		panic(err)
	}
	err = pg.Migrate()
	if err != nil {
		panic(err)
	}
}
