package data

type ICmdState interface {
	Exec(cmd string) error
	Stop() error
	Stack() *Stack
}