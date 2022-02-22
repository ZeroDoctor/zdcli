package comp

type exit int8

const (
	EXIT_SUC exit = iota
	EXIT_CMD
)

type ExitMessage struct {
	Code exit
	Msg  string
}

type Data struct {
	Type string
	Msg  interface{}
}
