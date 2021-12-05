package allone

import (
	"container/list"
	"math"
)

// Node define the node that
type Node struct {
	val int64
	set map[string]struct{}
}

func NewNode(val int64) *Node {
	return &Node{
		val: val,
		set: make(map[string]struct{}),
	}
}

func (n *Node) Erase(key string) {
	delete(n.set, key)
}

func (n *Node) Insert(key string) {
	n.set[key] = struct{}{}
}

func (n *Node) Has(key string) bool {
	_, ok := n.set[key]
	return ok
}

func (n *Node) Size() int64 {
	return int64(len(n.set))
}

func (n *Node) PickOneKey() string {
	for k := range n.set {
		return k
	}
	return ""
}

// -------------- AllOne --------------
type AllOne struct {
	data *list.List
	hash map[string]*list.Element
}

func Constructor() AllOne {
	data := list.New()
	data.PushFront(NewNode(math.MinInt64))
	data.PushBack(NewNode(math.MaxInt64))
	return AllOne{
		data: data,
		hash: make(map[string]*list.Element, 0),
	}
}

func (this *AllOne) add_to_right(ele *list.Element, key string, val int64) *list.Element {
	if ele.Next().Value.(*Node).val == val {
		ele.Next().Value.(*Node).Insert(key)
	} else {
		t := NewNode(val)
		t.Insert(key)
		this.data.InsertAfter(t, ele)
	}
	return ele.Next()
}

func (this *AllOne) add_to_left(ele *list.Element, key string, val int64) *list.Element {
	if ele.Prev().Value.(*Node).val == val {
		ele.Prev().Value.(*Node).Insert(key)
	} else {
		t := NewNode(val)
		t.Insert(key)
		this.data.InsertBefore(t, ele)
	}
	return ele.Prev()
}

func (this *AllOne) remove(node *list.Element) {
	this.data.Remove(node)
}

func (this *AllOne) Inc(key string) {
	if _, ok := this.hash[key]; !ok {
		this.hash[key] = this.add_to_right(this.data.Front(), key, 1)
	} else {
		ele := this.hash[key]
		node := ele.Value.(*Node)
		node.Erase(key)
		this.hash[key] = this.add_to_right(ele, key, node.val+1)
		if node.Size() == 0 {
			this.remove(ele)
		}
	}
}

func (this *AllOne) Dec(key string) {
	if _, ok := this.hash[key]; !ok {
		return
	}
	ele := this.hash[key]
	node := ele.Value.(*Node)
	node.Erase(key)
	if node.val > 1 {
		this.hash[key] = this.add_to_left(ele, key, node.val-1)
	} else {
		delete(this.hash, key)
	}
	if node.Size() == 0 {
		this.remove(ele)
	}
}

func (this *AllOne) GetMaxKey() string {
	if len(this.hash) > 0 {
		return this.data.Back().Prev().Value.(*Node).PickOneKey()
	}
	return ""
}

func (this *AllOne) GetMinKey() string {
	if len(this.hash) > 0 {
		return this.data.Front().Next().Value.(*Node).PickOneKey()
	}
	return ""
}

/**
 * Your AllOne object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Inc(key);
 * obj.Dec(key);
 * param_3 := obj.GetMaxKey();
 * param_4 := obj.GetMinKey();
 */
