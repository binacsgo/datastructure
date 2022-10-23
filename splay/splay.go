package splay

import (
	"strings"
)

type Comparable interface {
	// Compare defines a comparison function in splay that returns `true` if and only if the
	// current element is strictly greater than the incoming element.
	Compare(Comparable) bool
}

// StoredObj defines all the methods that need to be implemented by the element being stored.
type StoredObj interface {
	// Key returns the unique key used by the object in the splay.
	Key() string
	// String implements the String interface.
	String() string
	// Maintain defines the maintenance operation in the splay, which contains the properties
	// of the subtree rooted at the current node. We will update the properties of the current
	// node based on its left and right children.
	Maintain(StoredObj, StoredObj)

	Comparable
}

// storedObjForLookup defines one of the simplest StoredObj implementations for lookups only.
type storedObjForLookup struct{ key string }

func (o *storedObjForLookup) Key() string             { return o.key }
func (o *storedObjForLookup) String() string          { return o.key }
func (o *storedObjForLookup) Maintain(_, _ StoredObj) {}
func (o *storedObjForLookup) Compare(Comparable) bool { return false }
func NewStoredObjForLookup(key string) StoredObj {
	return &storedObjForLookup{
		key: key,
	}
}

var (
	_ Comparable = &storedObjForLookup{}
	_ StoredObj  = &storedObjForLookup{}
	_ Splay      = &splay{}

	minObj = NewStoredObjForLookup("MinObj")
	maxObj = NewStoredObjForLookup("MaxObj")
)

type (
	// RangeFunc visit objects by inorder traversal.
	RangeFunc func(StoredObj)
	// ConditionRangeFunc visit objects by inorder traversal.
	ConditionRangeFunc func(StoredObj) bool
)

type node struct {
	child  []*node
	parent *node
	obj    StoredObj
}

func newNode(o StoredObj, p *node) *node {
	return &node{
		child:  make([]*node, 2),
		parent: p,
		obj:    o,
	}
}

// getChildIndex indicates whether `x` is the right child of `y`.
func getChildIndex(x, y *node) int {
	if y != nil && y.child[1] == x {
		return 1
	}
	return 0
}

// Splay defines all methods of the splay-tree.
type Splay interface {
	// Insert a StoredObj into the splay. Returns true if successful.
	Insert(StoredObj) bool
	// Delete a StoredObj from the splay. Returns true if successful.
	Delete(StoredObj) bool
	// Get a StoredObj from the splay.
	Get(StoredObj) StoredObj
	// Partition will bring together all objects strictly smaller than the current object
	// in a subtree and return the root of the subtree.
	Partition(Comparable) StoredObj
	// Range traverses the entire splay in mid-order.
	Range(RangeFunc)
	// ConditionRange traverses the entire splay in mid-order and ends the access immediately
	// if ConditionRangeFunc returns false.
	ConditionRange(ConditionRangeFunc)
	// Len returns the number of all objects in the splay.
	Len() int
	// String implements the String interface.
	String() string
	// PrintTree outputs splay in the form of a tree diagram.
	PrintTree() string
}

type splay struct {
	root             *node
	minv, maxv       *node
	index            map[string]*node
	chooseChildIndex func(Comparable, *node) int
	maintain         func(*node)
}

func NewSplay() Splay {
	s := &splay{
		minv:  newNode(minObj, nil),
		maxv:  newNode(maxObj, nil),
		index: make(map[string]*node),
	}
	s.minv.child[1], s.maxv.parent = s.maxv, s.minv
	s.root = s.minv
	s.chooseChildIndex = func(o Comparable, n *node) int {
		if n == s.minv || n != s.maxv && o.Compare(n.obj) {
			return 1
		}
		return 0
	}
	s.maintain = func(n *node) {
		var leftChildObj, rightChildObj StoredObj
		if n.child[0] != nil && n.child[0] != s.minv {
			leftChildObj = n.child[0].obj
		}
		if n.child[1] != nil && n.child[1] != s.maxv {
			rightChildObj = n.child[1].obj
		}
		n.obj.Maintain(leftChildObj, rightChildObj)
	}
	return s
}

func (s *splay) Insert(v StoredObj) bool {
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

func (s *splay) Delete(v StoredObj) bool {
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

func (s *splay) Get(obj StoredObj) StoredObj {
	n, ok := s.index[obj.Key()]
	if !ok {
		return nil
	}
	return n.obj
}

func (s *splay) Partition(obj Comparable) StoredObj {
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

func (s *splay) Range(f RangeFunc) {
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

func (s *splay) ConditionRange(f ConditionRangeFunc) {
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

func (s *splay) Len() int {
	return len(s.index)
}

func (s *splay) String() string {
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

func (s *splay) PrintTree() string {
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
		output.WriteString(n.obj.String())
		output.WriteByte('\n')
		handleChild(n.child[0], true)
	}
	output.WriteString("SplayRoot\n")
	dfs(s.root, &strings.Builder{}, true)
	return output.String()
}

func (s *splay) rotate(x *node) {
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

func (s *splay) splay(x, k *node) {
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
