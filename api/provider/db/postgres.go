package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresDBReader struct {
	rows *sql.Rows
}

func (reader *PostgresDBReader) NextRow(dest ...interface{}) (bool, error) {
	if reader.rows.Next() {
		err := reader.rows.Scan(dest...)
		return err == nil, err
	}
	return false, reader.rows.Err()
}

func (reader *PostgresDBReader) GetRow(dest ...interface{}) error {
	found, err := reader.NextRow(dest...)
	if !found && err == nil {
		err = sql.ErrNoRows
	}
	return err
}

func (reader *PostgresDBReader) Close() {
	if reader.rows != nil {
		reader.rows.Close()
	}
}

type PostgresDBProvider struct {
	connection *sql.DB
}

func (provider *PostgresDBProvider) Query(query string, args ...interface{}) (DBReader, error) {
	response, err := provider.connection.Query(query, args...)
	return &PostgresDBReader{rows: response}, err
}

func (provider *PostgresDBProvider) Exec(query string, args ...interface{}) (int64, error) {
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

func (provider *PostgresDBProvider) Migrate() (uint, error) {
	driver, err := postgres.WithInstance(provider.connection, &postgres.Config{})
	if err != nil {
		return 0, err
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file:///migrations/postgres",
		"postgres", driver)
	if err != nil {
		return 0, err
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return 0, err
	}

	version, _, err := migrator.Version()
	if err != nil {
		return 0, err
	}

	return version, nil
}

func CreatePostgresDBProvider(URL string) (DBProvider, error) {
	db, err := sql.Open("postgres", URL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	provider := PostgresDBProvider{
		connection: db,
	}

	return &provider, nil
}
