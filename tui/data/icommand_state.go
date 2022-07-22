package data

type ICmdState interface {
	Exec(cmd string) error
	Stop() error
}

type ICmdStateManager interface {
	Exec(cmd string) error
	Stop() error
	SetStack(*Stack)
}
