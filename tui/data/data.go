package data

type exit int8

const (
	EXIT_SUC exit = iota
	EXIT_CMD
	EXIT_LUA
	EXIT_EDT
)

type ExitMessage struct {
	Code exit
	Msg  string
}

type Data struct {
	Type string
	Msg  interface{}
}
