// Package database implements the logic to connect to a database and run queries.
package database

import (
	"database/sql"
	"fmt"

	// Postgres driver to connect to the database.
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database type to connect to the database.
type Database struct {
	session      *gorm.DB
	conn         *sql.DB
	user         string
	password     string
	port         string
	databaseName string
}

// New returns a new database object.
func New(user, password, dbname, port string) *Database {
	return &Database{
		user:         user,
		password:     password,
		port:         port,
		databaseName: dbname,
	}
}

func (db *Database) newSession() error {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", db.user,
		db.password, db.port, db.databaseName)
	session, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("unable to start the database session: %w", err)
	}

	db.session = session

	return nil
}

// Connect to the database.
func (db *Database) Connect() error {
	if db.session == nil {
		if err := db.newSession(); err != nil {
			return err
		}
	}

	sqlDB, err := db.session.DB()
	if err != nil {
		return fmt.Errorf("unable to get the database connection: %w", err)
	}

	db.conn = sqlDB

	return nil
}

// Disconnect from the database.
func (db *Database) Disconnect() error {
	if db == nil {
		return nil
	}

	if db.conn == nil {
		return nil
	}

	if err := db.conn.Close(); err != nil {
		return err
	}

	db.conn = nil
	db.session = nil

	return nil
}
