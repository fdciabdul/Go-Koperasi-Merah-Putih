package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"koperasi-merah-putih/config"
)

type CassandraDB struct {
	Session *gocql.Session
}

func NewCassandraConnection(cfg *config.CassandraConfig) (*CassandraDB, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = cfg.Keyspace
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second

	if cfg.Username != "" && cfg.Password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	switch cfg.Consistency {
	case "one":
		cluster.Consistency = gocql.One
	case "quorum":
		cluster.Consistency = gocql.Quorum
	case "all":
		cluster.Consistency = gocql.All
	default:
		cluster.Consistency = gocql.Quorum
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Cassandra: %v", err)
	}

	log.Println("Connected to Cassandra successfully")
	return &CassandraDB{Session: session}, nil
}

func (c *CassandraDB) Close() {
	if c.Session != nil {
		c.Session.Close()
	}
}