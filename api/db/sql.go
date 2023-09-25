package db

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type SQLDBReader struct {
	rows *sql.Rows
}

func (reader *SQLDBReader) NextRow(dest ...interface{}) (bool, error) {
	if reader.rows.Next() {
		err := reader.rows.Scan(dest...)
		return err == nil, err
	}
	return false, reader.rows.Err()
}

func (reader *SQLDBReader) GetRow(dest ...interface{}) error {
	found, err := reader.NextRow(dest...)
	if !found && err == nil {
		err = sql.ErrNoRows
	}
	return err
}

func (reader *SQLDBReader) Close() {
	if reader.rows != nil {
		reader.rows.Close()
	}
}

type SQLDBProvider struct {
	connection *sql.DB
}

func (provider *SQLDBProvider) Query(query string, args ...interface{}) (DBReader, error) {
	response, err := provider.connection.Query(query, args...)
	return &SQLDBReader{rows: response}, err
}

func (provider *SQLDBProvider) Exec(query string, args ...interface{}) (int64, error) {
	result, err := provider.connection.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (provider *SQLDBProvider) Migrate() error {
	driver, err := postgres.WithInstance(provider.connection, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if err == migrate.ErrNoChange {
		log.Println("No changes applied to database")
	} else if err != nil {
		return err
	}

	version, dirty, err := migrator.Version()
	if err != nil {
		return err
	}

	log.Printf("Database version: %v, dirty - %v", version, dirty)
	return nil
}

func CreatePostgreSQLDBProvider(URL string) (DBProvider, error) {
	db, err := sql.Open("postgres", URL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	provider := SQLDBProvider{
		connection: db,
	}

	return &provider, nil
}
