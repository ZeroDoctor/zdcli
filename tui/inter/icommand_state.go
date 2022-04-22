package inter

type ICommandState interface {
	Exec(cmd string) error
	Stop() error
}
