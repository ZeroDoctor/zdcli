package data

import (
	"fmt"
	"sync"
)

type Stack struct {
	slice []interface{}
	m     sync.Mutex
}

func NewStack(slice ...interface{}) Stack {
	return Stack{slice: slice}
}

func (s *Stack) Clear() {
	s.m.Lock()
	defer s.m.Unlock()

	s.slice = []interface{}{}
}

func (s *Stack) Peek() interface{} {
	s.m.Lock()
	defer s.m.Unlock()

	return s.slice[0]
}

func (s *Stack) Push(elem interface{}) {
	s.m.Lock()
	defer s.m.Unlock()

	var arr []interface{}
	arr = append(arr, elem)
	arr = append(arr, s.slice...)
	s.slice = arr
}

func (s *Stack) Pop() interface{} {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.slice) <= 0 {
		return nil
	}

	elem := s.slice[0]
	s.slice = s.slice[1:]

	return elem
}

func (s *Stack) String() string {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.slice) <= 0 {
		return "[]"
	}

	var str string
	str = "["
	for _, s := range s.slice {
		str += fmt.Sprint(s, ",")
	}
	str = str[:len(str)-1]
	str += "]"
	return str
}

func (s *Stack) Len() int {
	s.m.Lock()
	defer s.m.Unlock()

	return len(s.slice)
}
