package allone

import (
	"container/list"
	"math"
)

// Node define the node that
type Node struct {
	val int64
	set map[any]struct{}
}

func NewNode(val int64) *Node {
	return &Node{
		val: val,
		set: make(map[any]struct{}),
	}
}

func (n *Node) Erase(key any) {
	delete(n.set, key)
}

func (n *Node) Insert(key any) {
	n.set[key] = struct{}{}
}

func (n *Node) Has(key any) bool {
	_, ok := n.set[key]
	return ok
}

func (n *Node) Size() int64 {
	return int64(len(n.set))
}

func (n *Node) PickOneKey() any {
	for k := range n.set {
		return k
	}
	return ""
}

// -------------- AllOneMinMax --------------
type AllOneMinMax struct {
	data *list.List
	hash map[any]*list.Element
}

func AllOneMinMaxConstructor() AllOneMinMax {
	data := list.New()
	data.PushFront(NewNode(math.MinInt64))
	data.PushBack(NewNode(math.MaxInt64))
	return AllOneMinMax{
		data: data,
		hash: make(map[any]*list.Element, 0),
	}
}

func (this *AllOneMinMax) add_to_right(ele *list.Element, key any, val int64) *list.Element {
	if ele.Next().Value.(*Node).val == val {
		ele.Next().Value.(*Node).Insert(key)
	} else {
		t := NewNode(val)
		t.Insert(key)
		this.data.InsertAfter(t, ele)
	}
	return ele.Next()
}

func (this *AllOneMinMax) add_to_left(ele *list.Element, key any, val int64) *list.Element {
	if ele.Prev().Value.(*Node).val == val {
		ele.Prev().Value.(*Node).Insert(key)
	} else {
		t := NewNode(val)
		t.Insert(key)
		this.data.InsertBefore(t, ele)
	}
	return ele.Prev()
}

func (this *AllOneMinMax) remove(node *list.Element) {
	this.data.Remove(node)
}

func (this *AllOneMinMax) Inc(key any) {
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

func (this *AllOneMinMax) Dec(key any) {
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

func (this *AllOneMinMax) GetMaxKey() any {
	if len(this.hash) > 0 {
		return this.data.Back().Prev().Value.(*Node).PickOneKey()
	}
	return ""
}

func (this *AllOneMinMax) GetMinKey() any {
	if len(this.hash) > 0 {
		return this.data.Front().Next().Value.(*Node).PickOneKey()
	}
	return ""
}

/**
 * Your AllOneMinMax object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Inc(key);
 * obj.Dec(key);
 * param_3 := obj.GetMaxKey();
 * param_4 := obj.GetMinKey();
 */
