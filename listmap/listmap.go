package listmap

type StoredObj interface {
	// GetGeneration() int64
}

type ProcessFunc func(StoredObj)

type ListItem struct {
	StoredObj
	next *ListItem
	prev *ListItem
}

func newListItem(obj StoredObj) *ListItem {
	return &ListItem{
		StoredObj: obj,
	}
}

// Next return the next item's pointer.
func (item *ListItem) Next() *ListItem {
	if item == nil {
		return nil
	}
	return item.next
}

type ListMap struct {
	head *ListItem
	hash map[string]*ListItem
}

func NewListMap() *ListMap {
	return &ListMap{
		hash: make(map[string]*ListItem),
	}
}

func (m *ListMap) Get(key string) StoredObj {
	item, ok := m.hash[key]
	if !ok || item == nil {
		return nil
	}
	return item.StoredObj
}

func (m *ListMap) Set(key string, obj StoredObj) {
	var item *ListItem
	if item = m.hash[key]; item != nil {
		m.remove(item)
		item.StoredObj = obj
	} else {
		item = newListItem(obj)
	}
	if m.head != nil {
		m.head.prev = item
	}
	item.next, item.prev = m.head, nil
	m.head = item
	m.hash[key] = item
}

func (m *ListMap) Delete(key string) {
	if item, ok := m.hash[key]; ok && item != nil {
		m.remove(item)
	}
	delete(m.hash, key)
}

func (m *ListMap) Front() *ListItem {
	return m.head
}

func (m *ListMap) Range(generation int64, processFunc ProcessFunc) {
	for e := m.head; e != nil; e = e.Next() {
		// if e.GetGeneration() <= generation {
		// 	break
		// }
		processFunc(e.StoredObj)
	}
}

func (m *ListMap) remove(item *ListItem) {
	if item == nil {
		return
	}
	if item.prev != nil {
		item.prev.next = item.next
	}
	if item.next != nil {
		item.next.prev = item.prev
	}
	if m.head == item {
		m.head = item.next
	}
}
