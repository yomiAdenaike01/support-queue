package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yomiAdenaike01/support-queue/internals/config"
	"github.com/yomiAdenaike01/support-queue/internals/db/queries"
)

const defaultURL = "host=localhost port=5432 user=user password=pwd dbname=supportops sslmode=disable"

func NewClient(config *config.Config) *sqlx.DB {
	databaseURL := config.GetEnvOrDefault("DATABASE_URL", defaultURL)
	db := sqlx.MustConnect("postgres", databaseURL)
	if db == nil {
		panic("Failed to connect to db")
	}

	db.MustExec(queries.SCHEMA_SQL_QUERY)
	return db
}
