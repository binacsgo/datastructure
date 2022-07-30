package set

import "fmt"

type Set interface {
	Exist(item any) bool
	Insert(item ...any)
	Delete(item ...any)
	Len() int
	List() []any
}

type SetImpl struct {
	data map[any]struct{}
}

func NewSet() Set {
	return &SetImpl{
		data: make(map[any]struct{}),
	}
}

func (s *SetImpl) Exist(item any) bool {
	if s == nil {
		return false
	}
	_, ok := s.data[item]
	return ok
}

func (s *SetImpl) Insert(items ...any) {
	if s == nil {
		return
	}
	for _, item := range items {
		if _, ok := s.data[item]; ok {
			continue
		}
		s.data[item] = struct{}{}
	}
}

func (s *SetImpl) Delete(items ...any) {
	if s == nil {
		return
	}
	for _, item := range items {
		if _, ok := s.data[item]; !ok {
			continue
		}
		delete(s.data, item)
	}
}

func (s *SetImpl) Len() int {
	if s == nil {
		return 0
	}
	return len(s.data)
}

func (s *SetImpl) List() []any {
	if s == nil {
		return nil
	}
	ret := make([]any, 0)
	for k := range s.data {
		ret = append(ret, k)
	}
	return ret
}

func (s *SetImpl) String() string {
	if s == nil {
		return "{}"
	}
	var str string
	for k := range s.data {
		str += fmt.Sprintf("{%v},", s.data[k])
	}
	return str
}
