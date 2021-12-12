package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	log "github.com/mj23978/chat-backend-x/logger/zerolog"
	"github.com/pkg/errors"
	"github.com/scylladb/gocqlx/v2"
	"regexp"
	"strings"
	"time"
)

type Config struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	Keyspace string `mapstructure:"keyspace"`
}

type Cassandra struct {
	Cluster *gocql.ClusterConfig
	Session *gocql.Session
}

func NewCassandra(c Config) *Cassandra {
	cas := &Cassandra{}

	cluster := gocql.NewCluster(c.Address)
	cluster.Keyspace = c.Keyspace
	cluster.Port = 9042
	cluster.ProtoVersion = 4
	cluster.Consistency = gocql.Quorum
	cas.Cluster = cluster
	session, err := cluster.CreateSession()
	if err != nil {
		log.Errorf("Cant Create Session: %v", err)
	}
	cas.Session = session

	log.Infof("cassandra new client : %v", c)
	return cas
}

type Table struct {
	Name        string
	Type        interface{}
	PrimaryKeys [][]string
	Indexes     []string
}

func (t *Table) CreateTable(cassandra *Cassandra) error {

	var err error
	if cassandra.Session == nil {
		cassandra.Session, err = cassandra.Cluster.CreateSession()
		if err != nil {
			return errors.Wrap(err, "Error On Cassandra.create %s")
		}
	}

	solo := []string{}
	tuple := [][]string{}
	for _, value := range t.PrimaryKeys {
		if len(value) == 1 {
			solo = append(solo, value[0])
		} else {
			tuple = append(tuple, value)
		}
	}
	pkString := "PRIMARY KEY ("
	for _, value := range tuple {
		pkString = pkString + "(" + strings.Join(value, ", ") + "), "
	}
	for _, value := range solo {
		pkString = pkString + value
	}
	pkString = pkString + ")"

	fields := []string{}

	conv := SchemaConverter{Enums: map[string]gocql.Type{"StatusType": gocql.TypeInt}}
	m, ok := conv.StructToMap(t.Type)
	if !ok {
		panic("Unable to get map from struct during create table")
	}

	for key, value := range m {
		key = ToSnakeCase(key)
		typ, err := conv.stringTypeOf(value)
		if err != nil {
			return err
		}
		fields = append(fields, fmt.Sprintf(`%q %v`, key, typ))
	}

	// Add primary key value to fields list
	fields = append(fields, pkString)
	err = cassandra.Session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %q.%q (%s)`, cassandra.Cluster.Keyspace, t.Name, strings.Join(fields, ", "))).Exec()
	if err != nil {
		return errors.Wrap(err, "Error On Cassandra.create %s")
	}
	for _, value := range t.Indexes {
		err = cassandra.Session.Query(fmt.Sprintf(`CREATE INDEX ON %q.%q (%s)`, cassandra.Cluster.Keyspace, t.Name, value)).Exec()
		if err != nil {
			return errors.Wrap(err, "Error On Cassandra.create %s")
		}
	}
	return nil

}

// CreateKeyspace creates keyspace with SimpleStrategy and RF derived from flags.
func (c *Cassandra) CreateKeyspace(ks string) error {
	c.Cluster.Timeout = 30 * time.Second

	session, err := gocqlx.WrapSession(c.Cluster.CreateSession())
	if err != nil {
		return err
	}
	defer session.Close()

	err = session.ExecStmt(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class' : 'SimpleStrategy', 'replication_factor' : 1}`, ks))
	if err != nil {
		return fmt.Errorf("create keyspace: %w", err)
	}

	return nil
}

func (c *Cassandra) CreateType(name string, typ interface{}) error {
	c.Cluster.Timeout = 30 * time.Second

	session, err := gocqlx.WrapSession(c.Cluster.CreateSession())
	if err != nil {
		return err
	}
	defer session.Close()

	conv := SchemaConverter{Enums: map[string]gocql.Type{"StatusType": gocql.TypeInt}}
	m, ok := conv.StructToMap(typ)
	if !ok {
		panic("Unable to get map from struct during create table")
	}

	fields := []string{}

	for key, value := range m {
		key = ToSnakeCase(key)
		typ, err := conv.stringTypeOf(value)
		if err != nil {
			return err
		}
		fields = append(fields, fmt.Sprintf(`%q %v`, key, typ))
	}

	err = session.ExecStmt(fmt.Sprintf(`CREATE TYPE IF NOT EXISTS %v (%s)`, name, strings.Join(fields, ", ")))
	if err != nil {
		return fmt.Errorf("create keyspace: %w", err)
	}

	return nil
}

func (c *Cassandra) CreateEnum(name string, typ interface{}) error {
	c.Cluster.Timeout = 30 * time.Second

	session, err := gocqlx.WrapSession(c.Cluster.CreateSession())
	if err != nil {
		return err
	}
	defer session.Close()
	conv := SchemaConverter{Enums: map[string]gocql.Type{"StatusType": gocql.TypeInt}}

	m, ok := conv.StructToMap(typ)
	if !ok {
		panic("Unable to get map from struct during create table")
	}

	fields := []string{}

	for key, value := range m {
		key = ToSnakeCase(key)
		typ, err := conv.stringTypeOf(value)
		if err != nil {
			return err
		}
		fields = append(fields, fmt.Sprintf(`%q %v`, key, typ))
	}

	err = session.ExecStmt(fmt.Sprintf(`CREATE TYPE IF NOT EXISTS %v (%s)`, name, strings.Join(fields, ", ")))
	if err != nil {
		return fmt.Errorf("create keyspace: %w", err)
	}

	return nil
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
