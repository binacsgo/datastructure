package dynamic

import (
	"strings"

	"github.com/binacsgo/datastructure/splay"
)

var (
	_ splay.Splay = &dynamicSplay{}
)

type node struct {
	child  []*node
	parent *node
	obj    splay.StoredObj
	info   splay.MaintainInfo
}

func newNode(o splay.StoredObj, p *node) *node {
	return &node{
		child:  make([]*node, 2),
		parent: p,
		obj:    o,
		info:   o.MakeMaintainInfo(),
	}
}

type dynamicSplay struct {
	root             *node
	minv, maxv       *node
	index            map[string]*node
	chooseChildIndex func(splay.Comparable, *node) int
	maintain         func(*node)
}

func New() splay.Splay {
	s := &dynamicSplay{
		minv:  newNode(splay.MinObj, nil),
		maxv:  newNode(splay.MaxObj, nil),
		index: make(map[string]*node),
	}
	s.minv.child[1], s.maxv.parent = s.maxv, s.minv
	s.root = s.minv
	s.chooseChildIndex = func(o splay.Comparable, n *node) int {
		if n == s.minv || n != s.maxv && o.Compare(n.obj) {
			return 1
		}
		return 0
	}
	s.maintain = func(n *node) {
		var leftChildInfo, rightChildInfo splay.MaintainInfo
		if n.child[0] != nil && n.child[0] != s.minv {
			leftChildInfo = n.child[0].info
		}
		if n.child[1] != nil && n.child[1] != s.maxv {
			rightChildInfo = n.child[1].info
		}
		n.info.Maintain(leftChildInfo, rightChildInfo)
	}
	return s
}

func (s *dynamicSplay) Insert(v splay.StoredObj) bool {
	if _, ok := s.index[v.Key()]; ok {
		return false
	}
	n := s.root
	var p *node
	for n != nil {
		p, n = n, n.child[s.chooseChildIndex(v, n)]
	}
	n = newNode(v, p)
	s.index[v.Key()] = n
	if p != nil {
		p.child[s.chooseChildIndex(v, p)] = n
	}
	s.splay(n, nil)
	return true
}

func (s *dynamicSplay) Delete(v splay.StoredObj) bool {
	n, ok := s.index[v.Key()]
	if !ok {
		return false
	}
	s.splay(n, nil)
	find := func(i int) (ret *node) {
		for ret = n.child[i]; ret.child[i^1] != nil; ret = ret.child[i^1] {
		}
		return
	}
	pre, nxt := find(0), find(1)
	s.splay(pre, nil)
	s.splay(nxt, pre)
	nxt.child[0] = nil
	s.maintain(nxt)
	s.maintain(pre)
	delete(s.index, v.Key())
	return true
}

func (s *dynamicSplay) Get(obj splay.StoredObj) splay.StoredObj {
	n, ok := s.index[obj.Key()]
	if !ok {
		return nil
	}
	return n.obj
}

func (s *dynamicSplay) Partition(obj splay.Comparable) splay.StoredObj {
	s.splay(s.minv, nil)
	var next *node
	for p := s.root; p != nil; {
		if s.chooseChildIndex(obj, p) == 1 {
			p = p.child[1]
		} else {
			next = p
			p = p.child[0]
		}
	}
	s.splay(next, s.minv)
	if next.child[0] == nil {
		return nil
	}
	return next.child[0].obj
}

func (s *dynamicSplay) Range(f splay.RangeFunc) {
	var dfs func(n *node)
	dfs = func(n *node) {
		if n == nil {
			return
		}
		dfs(n.child[0])
		if n != s.minv && n != s.maxv {
			f(n.obj)
		}
		dfs(n.child[1])
	}
	dfs(s.root)
}

func (s *dynamicSplay) ConditionRange(f splay.ConditionRangeFunc) {
	var dfs func(n *node)
	dfs = func(n *node) {
		if n == nil {
			return
		}
		dfs(n.child[0])
		if n != s.minv && n != s.maxv {
			if !f(n.obj) {
				return
			}
		}
		dfs(n.child[1])
	}
	dfs(s.root)
}

func (s *dynamicSplay) Len() int {
	return len(s.index)
}

func (s *dynamicSplay) String() string {
	output := &strings.Builder{}
	var dfs func(*node)
	dfs = func(n *node) {
		if n == nil {
			return
		}
		dfs(n.child[0])
		if n != s.minv && n != s.maxv {
			output.WriteString(n.obj.Key() + ",")
		}
		dfs(n.child[1])
	}
	dfs(s.root)
	return output.String()
}

func (s *dynamicSplay) Clone() splay.Splay {
	clone := New().(*dynamicSplay)
	index := make(map[string]*node, len(s.index))
	var dfs func(*node, *node) *node
	dfs = func(n, p *node) *node {
		if n == nil {
			return nil
		}
		var nn *node
		if n == s.minv {
			nn = clone.minv
			nn.parent = p
		} else if n == s.maxv {
			nn = clone.maxv
			nn.parent = p
		} else {
			nn = newNode(n.obj, p)
			nn.info = n.info.Clone() // ATTENTION
			index[nn.obj.Key()] = nn
		}
		nn.child[0], nn.child[1] = dfs(n.child[0], nn), dfs(n.child[1], nn)
		return nn
	}
	clone.root = dfs(s.root, nil)

	clone.index = index
	return clone
}

func (s *dynamicSplay) PrintTree() string {
	output := &strings.Builder{}
	var dfs func(*node, *strings.Builder, bool)
	dfs = func(n *node, prefixBuilder *strings.Builder, isBottom bool) {
		prefix := prefixBuilder.String()
		handleChild := func(n *node, flag bool) {
			if n == nil {
				return
			}
			nextPrefixBuilder := &strings.Builder{}
			nextPrefixBuilder.WriteString(prefix)
			if isBottom != flag {
				nextPrefixBuilder.WriteString("│   ")
			} else {
				nextPrefixBuilder.WriteString("    ")
			}
			dfs(n, nextPrefixBuilder, flag)
		}
		handleChild(n.child[1], false)
		output.WriteString(prefix)
		if isBottom {
			output.WriteString("└── ")
		} else {
			output.WriteString("┌── ")
		}
		output.WriteString(n.obj.String() + "[" + n.info.String() + "]")
		output.WriteByte('\n')
		handleChild(n.child[0], true)
	}
	output.WriteString("SplayRoot\n")
	dfs(s.root, &strings.Builder{}, true)
	return output.String()
}

// getChildIndex indicates whether `x` is the right child of `y`.
func getChildIndex(x, y *node) int {
	if y != nil && y.child[1] == x {
		return 1
	}
	return 0
}

func (s *dynamicSplay) rotate(x *node) {
	y := x.parent
	z := y.parent
	k := getChildIndex(x, y)
	if z != nil {
		z.child[getChildIndex(y, z)] = x
	}
	x.parent = z
	y.child[k] = x.child[k^1]
	if x.child[k^1] != nil {
		x.child[k^1].parent = y
	}
	x.child[k^1] = y
	y.parent = x
	s.maintain(y)
	s.maintain(x)
}

func (s *dynamicSplay) splay(x, k *node) {
	for x.parent != k {
		y := x.parent
		z := y.parent
		if z != k {
			if getChildIndex(x, y) != getChildIndex(y, z) {
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
