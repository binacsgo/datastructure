package static

import (
	"strconv"
	"strings"

	"github.com/binacsgo/datastructure/splay"
)

type node struct {
	lchild, rchild int
	parent         int
	key            string
	obj            splay.StoredObj
}

func newNode(o splay.StoredObj, p int) node {
	return node{
		parent: p,
		key:    o.Key(),
		obj:    o,
	}
}

type staticSplay struct {
	root       int
	minv, maxv int
	hash       map[string]int
	items      []node
	infos      []splay.MaintainInfo
	count      int

	chooseChildIndex func(splay.Comparable, int) int
	maintain         func(int)
}

func New() splay.Splay {
	s := &staticSplay{
		minv:  1,
		maxv:  2,
		hash:  make(map[string]int),
		items: []node{newNode(splay.NilObj, -1), newNode(splay.MinObj, 0), newNode(splay.MaxObj, 1)},
		infos: []splay.MaintainInfo{splay.NilObj.MakeMaintainInfo(), splay.MinObj.MakeMaintainInfo(), splay.MaxObj.MakeMaintainInfo()},
		count: 2,
	}
	s.items[s.minv].rchild, s.items[s.maxv].parent = s.maxv, s.minv
	s.root = s.minv
	s.chooseChildIndex = func(o splay.Comparable, n int) int {
		if n == s.minv || n != s.maxv && o.Compare(s.items[n].obj) {
			return 1
		}
		return 0
	}
	s.maintain = func(i int) {
		n := &s.items[i]
		var leftChildInfo, rightChildInfo splay.MaintainInfo
		if n.lchild != 0 && n.lchild != s.minv {
			leftChildInfo = s.infos[n.lchild]
		}
		if n.rchild != 0 && n.rchild != s.maxv {
			rightChildInfo = s.infos[n.rchild]
		}
		s.infos[i].Maintain(leftChildInfo, rightChildInfo)
	}
	return s
}

func (s *staticSplay) Insert(v splay.StoredObj) bool {
	if i, ok := s.hash[v.Key()]; ok {
		s.items[i].obj = v
		return false
	}
	p, n := 0, s.root
	for n != 0 {
		p = n
		if s.chooseChildIndex(v, n) == 1 {
			n = s.items[n].rchild
		} else {
			n = s.items[n].lchild
		}
	}

	{
		s.items = append(s.items, newNode(v, p))
		s.infos = append(s.infos, v.MakeMaintainInfo())
		s.count++
		s.hash[v.Key()] = s.count
		n = s.count
	}
	if p != 0 {
		if s.chooseChildIndex(v, p) == 1 {
			s.items[p].rchild = n
		} else {
			s.items[p].lchild = n
		}
	}
	s.splay(n, 0)
	return true
}

func (s *staticSplay) Delete(v splay.StoredObj) bool {
	i, ok := s.hash[v.Key()]
	if !ok {
		return false
	}
	s.splay(i, 0)
	find := func(dir int) (ret int) {
		if dir == 0 {
			for ret = s.items[i].lchild; s.items[ret].rchild != 0; ret = s.items[ret].rchild {
			}
		} else {
			for ret = s.items[i].rchild; s.items[ret].lchild != 0; ret = s.items[ret].lchild {
			}
		}
		return
	}
	pre, nxt := find(0), find(1)
	s.splay(pre, 0)
	s.splay(nxt, pre)
	s.items[nxt].lchild = 0
	s.maintain(nxt)
	s.maintain(pre)
	delete(s.hash, v.Key())

	lastIndex := s.count
	if i != lastIndex {
		lastNode := &s.items[lastIndex]
		lastParent, lastLeftChild, lastRightChild := lastNode.parent, lastNode.lchild, lastNode.rchild
		if lastParent != 0 {
			if s.getChildIndex(lastIndex, lastParent) == 1 {
				s.items[lastParent].rchild = i
			} else {
				s.items[lastParent].lchild = i
			}
		}
		if lastLeftChild != 0 {
			s.items[lastLeftChild].parent = i
		}
		if lastRightChild != 0 {
			s.items[lastRightChild].parent = i
		}
		s.hash[lastNode.key] = i
		s.items[i] = s.items[lastIndex]
		s.infos[i] = s.infos[lastIndex]
		if s.root == lastIndex {
			s.root = i
		}
	}
	s.items = s.items[:s.count]
	s.infos = s.infos[:s.count]
	s.count--
	return true
}

func (s *staticSplay) Get(obj splay.StoredObj) splay.StoredObj {
	i, ok := s.hash[obj.Key()]
	if !ok {
		return nil
	}
	return s.items[i].obj
}

func (s *staticSplay) Partition(obj splay.Comparable) splay.StoredObj {
	s.splay(s.minv, 0)
	var next int
	for p := s.root; p != 0; {
		if s.chooseChildIndex(obj, p) == 1 {
			p = s.items[p].rchild
		} else {
			next = p
			p = s.items[p].lchild
		}
	}
	s.splay(next, s.minv)
	p := s.items[next].lchild
	if p == 0 {
		return nil
	}
	return s.items[p].obj
}

func (s *staticSplay) Range(f splay.RangeFunc) {
	var dfs func(int)
	dfs = func(i int) {
		if i == 0 {
			return
		}
		dfs(s.items[i].lchild)
		if i != s.minv && i != s.maxv {
			f(s.items[i].obj)
		}
		dfs(s.items[i].rchild)
	}
	dfs(s.root)
}

func (s *staticSplay) ConditionRange(f splay.ConditionRangeFunc) {
	var dfs func(int)
	dfs = func(i int) {
		if i == 0 {
			return
		}
		dfs(s.items[i].lchild)
		if i != s.minv && i != s.maxv {
			if !f(s.items[i].obj) {
				return
			}
		}
		dfs(s.items[i].rchild)
	}
	dfs(s.root)
}

func (s *staticSplay) Len() int {
	return s.count - 2
}

func (s *staticSplay) String() string {
	output := &strings.Builder{}
	var dfs func(int)
	dfs = func(i int) {
		if i == 0 {
			return
		}
		dfs(s.items[i].lchild)
		if i != s.minv && i != s.maxv {
			output.WriteString(s.items[i].key + ",")
		}
		dfs(s.items[i].rchild)
	}
	dfs(s.root)
	return output.String()
}

func (s *staticSplay) Clone() splay.Splay {
	clone := New().(*staticSplay)
	hash, items, infos := make(map[string]int, len(s.hash)), make([]node, len(s.hash)+3), make([]splay.MaintainInfo, len(s.hash)+3)

	copy(items, s.items)
	len := len(s.items)
	// TODO: Improve the `infos`.
	infos[0], infos[1], infos[2] = s.infos[0].Clone(), s.infos[1].Clone(), s.infos[2].Clone()
	for i := 3; i < len; i++ {
		hash[items[i].key] = i
		infos[i] = s.infos[i].Clone()
	}

	clone.hash, clone.items, clone.infos, clone.count, clone.root = hash, items, infos, s.count, s.root
	return clone
}

func (s *staticSplay) PrintTree() string {
	output := &strings.Builder{}
	var dfs func(int, *strings.Builder, bool)
	dfs = func(i int, prefixBuilder *strings.Builder, isBottom bool) {
		prefix := prefixBuilder.String()
		handleSon := func(j int, flag bool) {
			if j == 0 {
				return
			}
			nextPrefixBuilder := &strings.Builder{}
			nextPrefixBuilder.WriteString(prefix)
			if isBottom != flag {
				nextPrefixBuilder.WriteString("│   ")
			} else {
				nextPrefixBuilder.WriteString("    ")
			}
			dfs(j, nextPrefixBuilder, flag)
		}
		handleSon(s.items[i].rchild, false)
		output.WriteString(prefix)
		if isBottom {
			output.WriteString("└── ")
		} else {
			output.WriteString("┌── ")
		}
		output.WriteString(s.items[i].obj.String() + "(" + strconv.Itoa(i) + ")" + "[" + s.infos[i].String() + "]")
		output.WriteByte('\n')
		handleSon(s.items[i].lchild, true)
	}
	output.WriteString("SplayRoot:" + "root=" + strconv.Itoa(s.root) + "\n")
	dfs(s.root, &strings.Builder{}, true)
	return output.String()
}

// getChildIndex indicates whether `x` is the right child of `y`.
func (s *staticSplay) getChildIndex(x, y int) int {
	if y != 0 && s.items[y].rchild == x {
		return 1
	}
	return 0
}

func (s *staticSplay) rotate(x int) {
	y := s.items[x].parent
	z := s.items[y].parent
	k := s.getChildIndex(x, y)
	if z != 0 {
		if s.getChildIndex(y, z) == 1 {
			s.items[z].rchild = x
		} else {
			s.items[z].lchild = x
		}
	}
	s.items[x].parent = z
	if k == 1 {
		s.items[y].rchild = s.items[x].lchild
		if s.items[x].lchild != 0 {
			s.items[s.items[x].lchild].parent = y
		}
		s.items[x].lchild = y
	} else {
		s.items[y].lchild = s.items[x].rchild
		if s.items[x].rchild != 0 {
			s.items[s.items[x].rchild].parent = y
		}
		s.items[x].rchild = y
	}
	s.items[y].parent = x
	s.maintain(y)
	s.maintain(x)
}

func (s *staticSplay) splay(x, k int) {
	for s.items[x].parent != k {
		y := s.items[x].parent
		z := s.items[y].parent
		if z != k {
			if s.getChildIndex(x, y) != s.getChildIndex(y, z) {
				s.rotate(x)
			} else {
				s.rotate(y)
			}
		}
		s.rotate(x)
	}
	if k == 0 {
		s.root = x
	}
}
