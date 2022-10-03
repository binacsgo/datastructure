package splay

import (
	"strings"
)

type StoredObj interface {
	Key() string
}

type storedObjForLookup struct {
	key string
}

func (o *storedObjForLookup) Key() string {
	return o.key
}

func NewStoredObjForLookup(key string) StoredObj {
	return &storedObjForLookup{
		key: key,
	}
}

var (
	_ StoredObj = &storedObjForLookup{}

	MinObj = NewStoredObjForLookup("MinObj")
	MaxObj = NewStoredObjForLookup("MaxObj")
)

type (
	// CmpFunc return true when o1 > o2
	CmpFunc func(StoredObj, StoredObj) bool

	// MaintainFunc maintain the StoredObj by its sons.
	MaintainFunc func(StoredObj, StoredObj, StoredObj)

	// ConditionRangeFunc visit objects by inorder traversal.
	ConditionRangeFunc func(StoredObj) bool
)

type node struct {
	son    []*node
	parent *node
	obj    StoredObj
	size   int
}

func newNode(o StoredObj, p *node) *node {
	return &node{
		son:    make([]*node, 2),
		parent: p,
		obj:    o,
	}
}

// getSonIndex indicates whether `x` is the right son of `y`.
func getSonIndex(x, y *node) int {
	if y != nil && y.son[1] == x {
		return 1
	}
	return 0
}

type Splay interface {
	Insert(StoredObj) bool
	Delete(StoredObj) bool
	Get(StoredObj) StoredObj
	ConditionRange(ConditionRangeFunc)
	Len() int
	Equal(Splay) bool
	String() string
}

type splay struct {
	root           *node
	minv, maxv     *node
	chooseSonIndex func(StoredObj, *node) int
	maintain       func(*node)
	index          map[string]*node
}

var _ Splay = &splay{}

func NewSplay(cmp CmpFunc, maintain MaintainFunc) Splay {
	s := &splay{
		minv:  newNode(MinObj, nil),
		maxv:  newNode(MaxObj, nil),
		index: make(map[string]*node),
	}
	s.minv.son[1], s.maxv.parent = s.maxv, s.minv
	s.root = s.minv
	s.chooseSonIndex = func(o StoredObj, n *node) int {
		if n == s.minv || n != s.maxv && cmp(o, n.obj) {
			return 1
		}
		return 0
	}
	s.maintain = func(n *node) {
		if maintain != nil {
			var leftSonObj, rightSonObj StoredObj
			if n.son[0] != s.minv {
				leftSonObj = n.son[0].obj
			}
			if n.son[1] != s.maxv {
				rightSonObj = n.son[1].obj
			}
			maintain(n.obj, leftSonObj, rightSonObj)
		}
	}
	return s
}

func (s *splay) rotate(x *node) {
	y := x.parent
	z := y.parent
	k := getSonIndex(x, y)
	if z != nil {
		z.son[getSonIndex(y, z)] = x
	}
	x.parent = z
	y.son[k] = x.son[k^1]
	if x.son[k^1] != nil {
		x.son[k^1].parent = y
	}
	x.son[k^1] = y
	y.parent = x
	s.maintain(y)
	s.maintain(x)
}

func (s *splay) splay(x, k *node) {
	for x.parent != k {
		y := x.parent
		z := y.parent
		if z != k {
			if getSonIndex(x, y) != getSonIndex(y, z) {
				s.rotate(x)
			} else {
				s.rotate(y)
			}
		}
		s.rotate(x)
	}
	if k == nil {
		s.root = x
	}
}

func (s *splay) Insert(v StoredObj) bool {
	if _, ok := s.index[v.Key()]; ok {
		return false
	}
	n := s.root
	var p *node
	for n != nil {
		p, n = n, n.son[s.chooseSonIndex(v, n)]
	}
	n = newNode(v, p)
	s.index[v.Key()] = n
	if p != nil {
		p.son[s.chooseSonIndex(v, p)] = n
	}
	s.splay(n, nil)
	return true
}

func (s *splay) Delete(v StoredObj) bool {
	n, ok := s.index[v.Key()]
	if !ok {
		return false
	}
	s.splay(n, nil)
	find := func(i int) (ret *node) {
		ret = n.son[i]
		for {
			if ret.son[i^1] == nil {
				return
			}
			ret = ret.son[i^1]
		}
	}
	pre, nxt := find(0), find(1)
	s.splay(pre, nil)
	s.splay(nxt, pre)
	nxt.son[0] = nil
	delete(s.index, v.Key())
	return true
}

func (s *splay) Get(obj StoredObj) StoredObj {
	n, ok := s.index[obj.Key()]
	if !ok {
		return nil
	}
	return n.obj
}

func (s *splay) ConditionRange(f ConditionRangeFunc) {
	var dfs func(n *node)
	dfs = func(n *node) {
		if n == nil {
			return
		}
		dfs(n.son[0])
		if n != s.minv && n != s.maxv {
			if !f(n.obj) {
				return
			}
		}
		dfs(n.son[1])
	}
	dfs(s.root)
}

func (s *splay) Len() int {
	return len(s.index)
}

func (s *splay) Equal(t Splay) bool {
	// TODO: FIX 同优先级这里的顺序可能不一样
	return s.String() == t.String()
}

func (s *splay) String() string {
	output := &strings.Builder{}
	var dfs func(n *node)
	dfs = func(n *node) {
		if n == nil {
			return
		}
		dfs(n.son[0])
		if n != s.minv && n != s.maxv {
			output.WriteString(n.obj.Key() + ",")
		}
		dfs(n.son[1])
	}
	dfs(s.root)
	return output.String()
}
