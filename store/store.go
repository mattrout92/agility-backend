package store

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	dBConnect = os.Getenv("DBCONNECT")
)

var mtx sync.RWMutex

// Config represents the configuration required
type Config struct {
	Ctx               context.Context `json:"-"`
	MaxSQLConnections int             `json:"max_sql_connections"`
}

// DB represents a db with xray, sql and/or sqlx
type DB struct {
	sqlx *sqlx.DB
}

var maxConnectionsQuery = "show status where `variable_name` = ('Max_used_connections')"

var queries = make(map[string]string)

// Connect connects to the db with the relevant packages
func Connect(cfg *Config) *DB {
	database := &DB{}

	ctx := cfg.Ctx

	var maxConns int

	log.Println("connecting with sqlx driver")

	db, err := sqlx.ConnectContext(ctx, "mysql", dBConnect+"/Agility")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.PingContext(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if maxConns == 0 {
		var k, v string
		err = db.QueryRowContext(ctx, maxConnectionsQuery).Scan(&k, &v)

		maxConsPerInstance, _ := strconv.Atoi(v)
		fmt.Println("Max sql connections TOTAL", v)

		maxConns := int(math.Round(float64(maxConsPerInstance) / 20))
		fmt.Println("Max sql connections per instance", maxConns)
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(2)

	database.sqlx = db

	fmt.Println("DBs Connected mfstore")

	return database
}

// GetQuery returns a query string
func (db *DB) GetQuery(name string) string {
	mtx.Lock()
	defer mtx.Unlock()
	query, ok := queries[name]
	if !ok {
		f, err := os.Open("./queries/" + name + ".sql")
		if err != nil {
			log.Fatal(err) // exit program is query file is missing - critical error
		}
		defer f.Close()

		fbuf := new(bytes.Buffer)
		fbuf.ReadFrom(f)

		query = string(fbuf.Bytes()[:])

		queries[name] = query

	}
	return query
}

// SQLX returns the sqlx db
func (db *DB) SQLX() *sqlx.DB {
	return db.sqlx
}
