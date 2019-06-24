package main

import "errors"

var (
	EmptyErr = errors.New("empty stack")
)

type stack []uint16

func (s *stack) Push(v uint16) {
	*s = append(*s, v)
}

func (s *stack) Pop() (uint16, error) {
	l := len(*s)

	if l == 0 {
		return 0, EmptyErr
	}

	ret := (*s)[l-1]
	*s = (*s)[:l-1]
	return ret, nil
}
