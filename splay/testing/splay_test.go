package splaytest

import (
	"strconv"
	"testing"

	"github.com/binacsgo/datastructure/splay"
	"github.com/binacsgo/datastructure/splay/dynamic"
	"github.com/binacsgo/datastructure/splay/static"
)

type info struct {
	obj *obj
	sz  int
}

func (o *info) Maintain(ls, rs splay.MaintainInfo) {
	o.sz = 1
	if ls != nil {
		o.sz += ls.(*info).sz
	}
	if rs != nil {
		o.sz += rs.(*info).sz
	}
}

func (o *info) Clone() splay.MaintainInfo {
	return &info{obj: o.obj, sz: o.sz}
}

func (o *info) String() string {
	return strconv.Itoa(o.sz)
}

type obj struct {
	k int
	v int
}

var (
	_ splay.Comparable = &obj{}
	_ splay.StoredObj  = &obj{}
)

func makeObj(k, v int) splay.StoredObj {
	return &obj{k: k, v: v}
}

func (o *obj) Key() string                          { return strconv.Itoa(o.k) }
func (o *obj) String() string                       { return o.Key() }
func (o *obj) MakeMaintainInfo() splay.MaintainInfo { return &info{obj: o, sz: 1} }
func (o *obj) Compare(so splay.Comparable) bool     { return o.v > so.(*obj).v }

func Test_Splay(t *testing.T) {
	for _, s := range []splay.Splay{dynamic.New(), static.New()} {
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
				t.Errorf("There are %v items in splay, expect %v", s.Len(), i*3)
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
}

func Test_Rotate(t *testing.T) {
	for _, s := range []splay.Splay{dynamic.New(), static.New()} {
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
}

func Test_splay_Clone(t *testing.T) {
	for _, s := range []splay.Splay{dynamic.New(), static.New()} {
		for i := 1; i < 10; i++ {
			for j := 1; j < 4; j++ {
				s.Insert(makeObj(i*10+j, i))
			}
		}

		o := s.Clone()

		if o.Len() != s.Len() {
			t.Errorf("Expect length %v, got %v\n", s.Len(), o.Len())
		}

		s.Range(func(so splay.StoredObj) {
			origin := so.(*obj)
			if o.Get(so) != origin {
				t.Errorf("Expect object %v, got %v\n", origin, o.Get(so))
			}
		})

		t.Logf(s.PrintTree())
		t.Logf(o.PrintTree())

		o.Partition(makeObj(59, 5))

		t.Logf(s.PrintTree())
		t.Logf(o.PrintTree())
	}
}
