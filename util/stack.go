package util

import (
	"fmt"
	"sync"
)

type Stack[T any] struct {
	slice []T
	m     sync.Mutex
}

func NewStack[T any](slice ...T) Stack[T] {
	return Stack[T]{slice: slice}
}

func (s *Stack[T]) Clear() {
	s.m.Lock()
	defer s.m.Unlock()

	s.slice = []T{}
}

func (s *Stack[T]) Peek() T {
	s.m.Lock()
	defer s.m.Unlock()

	return s.slice[0]
}

func (s *Stack[T]) Push(elem ...T) {
	s.m.Lock()
	defer s.m.Unlock()

	var arr []T
	arr = append(arr, elem...)
	arr = append(arr, s.slice...)
	s.slice = arr
}

func (s *Stack[T]) Pop() *T {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.slice) <= 0 {
		return nil
	}

	elem := s.slice[0]
	s.slice = s.slice[1:]

	return &elem
}

func (s *Stack[T]) String() string {
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

func (s *Stack[T]) Len() int {
	s.m.Lock()
	defer s.m.Unlock()

	return len(s.slice)
}
