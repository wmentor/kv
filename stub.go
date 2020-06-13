package lbase

import (
	"bytes"
	"container/list"
)

type stub struct {
	lst *list.List
}

func newStub() *stub {
	return &stub{lst: list.New()}
}

func (s *stub) Close() {
	s.lst = list.New()
}

func (s *stub) Delete(key []byte) {

	if s == nil || len(key) == 0 {
		return
	}

	for e := s.lst.Front(); e != nil; e = e.Next() {
		v := e.Value.([][]byte)
		cmp := bytes.Compare(v[0], key)
		if cmp == 0 {
			s.lst.Remove(e)
			return
		}
	}
}

func (s *stub) Get(key []byte) []byte {

	if s == nil || len(key) == 0 {
		return nil
	}

	for e := s.lst.Front(); e != nil; e = e.Next() {
		v := e.Value.([][]byte)
		cmp := bytes.Compare(v[0], key)
		if cmp == 0 {
			return v[1]
		}
	}

	return nil
}

func (s *stub) Has(key []byte) bool {
	return len(s.Get(key)) > 0
}

func (s *stub) Set(key []byte, value []byte) {

	if s == nil || len(key) == 0 {
		return
	}

	if len(value) == 0 {
		s.Delete(key)
		return
	}

	for e := s.lst.Front(); e != nil; e = e.Next() {
		v := e.Value.([][]byte)
		cmp := bytes.Compare(v[0], key)
		if cmp == 0 {
			v[1] = value
			return
		} else if cmp > 0 {
			s.lst.InsertBefore([][]byte{key, value}, e)
			return
		}
	}

	s.lst.PushBack([][]byte{key, value})
}
