package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	"cal.ope.cool/pkg/db/orm"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDB struct {
	db  *sql.DB
	Orm *orm.Queries
}

type SqliteOpts struct {
	ConnString string
	DbName     string
}

//go:embed sqlc/schema.sql
var ddl string

func New(opts *SqliteOpts) *SqliteDB {
	ctx := context.Background()

	if opts.ConnString == "" {
		opts.ConnString = os.Getenv("DB_URL")
	}
	if opts.DbName == "" {
		opts.DbName = "database"
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", opts.ConnString))
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}

	orm := orm.New(db)
	return &SqliteDB{db: db, Orm: orm}
}
