package view

type Data struct {
	Type string
	Msg  interface{}
}

func NewData(t, m string) Data {
	return Data{Type: t, Msg: m}
}
