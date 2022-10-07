package splay

import (
	"fmt"
	"strconv"
	"testing"
)

type obj struct {
	k  int
	v  int
	sz int
}

var (
	_ Comparable = &obj{}
	_ StoredObj  = &obj{}
)

func makeObj(k, v int) StoredObj {
	return &obj{k: k, v: v, sz: 1}
}

func (o *obj) Key() string                { return strconv.Itoa(o.k) }
func (o *obj) String() string             { return o.Key() + fmt.Sprintf("(%d)", o.sz) }
func (o *obj) Compare(so Comparable) bool { return o.v > so.(*obj).v }
func (o *obj) Maintain(ls, rs StoredObj) {
	o.sz = 1
	if ls != nil {
		o.sz += ls.(*obj).sz
	}
	if rs != nil {
		o.sz += rs.(*obj).sz
	}
}

func Test_Splay(t *testing.T) {
	s := NewSplay()
	t.Log(s.PrintTree())

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
		t.Logf("After i=%v got splay: %s\n", i, s)
	}
	for j := 1; j < 4; j++ {
		for i := 1; i < 10; i++ {
			if s.Get(makeObj(i*10+j, i)) == nil {
				t.Errorf("There should be %v in splay.", i*10+j)
			}
			s.Delete(makeObj(i*10+j, i))
		}
		t.Logf("After j=%v got splay: %s\n", j, s)
	}
}

func Test_Rotate(t *testing.T) {
	s := NewSplay()
	t.Log(s.PrintTree())

	for i := 1; i < 10; i++ {
		for j := 1; j < 4; j++ {
			s.Insert(makeObj(i*10+j, i))
		}
	}
	t.Log(s.PrintTree())
	for j := 1; j < 2; j++ {
		for i := 1; i < 10; i++ {
			if s.Get(makeObj(i*10+j, i)) == nil {
				t.Errorf("There should be %v in splay.", i*10+j)
			}
			s.Delete(makeObj(i*10+j, i))
		}
		t.Logf("After j=%v got splay: %s\n", j, s)
	}
	t.Log(s.PrintTree())
	s.Partition(makeObj(59, 5))
	t.Log(s.PrintTree())
	for j := 2; j < 4; j++ {
		for i := 1; i < 10; i++ {
			if s.Get(makeObj(i*10+j, i)) == nil {
				t.Errorf("There should be %v in splay.", i*10+j)
			}
			s.Delete(makeObj(i*10+j, i))
		}
		t.Logf("After j=%v got splay: %s\n", j, s)
	}
	t.Log(s.PrintTree())

	s.Partition(makeObj(59, 5))
}
