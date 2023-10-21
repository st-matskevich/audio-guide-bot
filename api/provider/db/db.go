package db

type DBReader interface {
	NextRow(dest ...interface{}) (bool, error)
	GetRow(dest ...interface{}) error
	Close()
}

type DBProvider interface {
	Query(query string, args ...interface{}) (DBReader, error)
	Exec(query string, args ...interface{}) (int64, error)
	Migrate() error
}
