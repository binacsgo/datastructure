package splay

import (
	"fmt"
	"strconv"
	"testing"
)

type obj struct {
	k int
	v int
}

var _ StoredObj = &obj{}

func makeObj(k, v int) StoredObj {
	return &obj{k: k, v: v}
}

func (o *obj) Key() string {
	if o == nil {
		return "empty"
	}
	return strconv.Itoa(o.k)
}

func Test_Splay(t *testing.T) {
	s := NewSplay(
		func(so1, so2 StoredObj) bool {
			o1, o2 := so1.(*obj), so2.(*obj)
			return o1.v > o2.v
		},
		nil,
	)

	for i := 1; i < 10; i++ {
		for j := 1; j < 4; j++ {
			if s.Get(makeObj(i*10+j, i)) != nil {
				t.Errorf("There shouldn't be %v in splay.", i*10+j)
			}
		}
		for j := 1; j < 4; j++ {
			s.Insert(makeObj(i*10+j, i))
		}
		if s.Len() != i*3 {
			t.Errorf("There are %v items in splay, expect %v", s.Len(), i*4)
		}
		fmt.Printf("After i=%v got splay: %s\n", i, s)
	}
	for j := 1; j < 4; j++ {
		for i := 1; i < 10; i++ {
			if s.Get(makeObj(i*10+j, i)) == nil {
				t.Errorf("There should be %v in splay.", i*10+j)
			}
			s.Delete(makeObj(i*10+j, i))
		}
		fmt.Printf("After j=%v got splay: %s\n", j, s)
	}

}
