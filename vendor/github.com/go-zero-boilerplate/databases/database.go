package databases

type Database interface {
	IsErrNoRows(err error) bool

	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (ExecResult, error)
	QueryRow(query string, args ...interface{}) ResultRow
	Query(query string, args ...interface{}) (ResultRows, error)
}
