package allone

import (
	"container/list"
	"math"
)

// Node define the node that
type Node struct {
	val int64
	set map[interface{}]struct{}
}

func NewNode(val int64) *Node {
	return &Node{
		val: val,
		set: make(map[interface{}]struct{}),
	}
}

func (n *Node) Erase(key interface{}) {
	delete(n.set, key)
}

func (n *Node) Insert(key interface{}) {
	n.set[key] = struct{}{}
}

func (n *Node) Has(key interface{}) bool {
	_, ok := n.set[key]
	return ok
}

func (n *Node) Size() int64 {
	return int64(len(n.set))
}

func (n *Node) PickOneKey() interface{} {
	for k := range n.set {
		return k
	}
	return ""
}

// -------------- AllOne --------------
type AllOne struct {
	data *list.List
	hash map[interface{}]*list.Element
}

func Constructor() AllOne {
	data := list.New()
	data.PushFront(NewNode(math.MinInt64))
	data.PushBack(NewNode(math.MaxInt64))
	return AllOne{
		data: data,
		hash: make(map[interface{}]*list.Element, 0),
	}
}

func (this *AllOne) add_to_right(ele *list.Element, key interface{}, val int64) *list.Element {
	if ele.Next().Value.(*Node).val == val {
		ele.Next().Value.(*Node).Insert(key)
	} else {
		t := NewNode(val)
		t.Insert(key)
		this.data.InsertAfter(t, ele)
	}
	return ele.Next()
}

func (this *AllOne) add_to_left(ele *list.Element, key interface{}, val int64) *list.Element {
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

func (this *AllOne) Inc(key interface{}) {
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

func (this *AllOne) Dec(key interface{}) {
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

func (this *AllOne) GetMaxKey() interface{} {
	if len(this.hash) > 0 {
		return this.data.Back().Prev().Value.(*Node).PickOneKey()
	}
	return ""
}

func (this *AllOne) GetMinKey() interface{} {
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
