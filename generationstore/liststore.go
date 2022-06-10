package generationstore

import (
	"fmt"
	"sort"
	"sync/atomic"
)

func nextGeneration(generation *uint64) uint64 {
	return atomic.AddUint64(generation, 1)
}

// ------------------- ListItem -------------------
type ListItem struct {
	key string
	StoredObj
	next *ListItem
	prev *ListItem
}

func newListItem(key string, obj StoredObj) *ListItem {
	return &ListItem{
		key:       key,
		StoredObj: obj,
	}
}

func (item *ListItem) Obj() StoredObj {
	if item == nil {
		return nil
	}
	return item.StoredObj
}

func (item *ListItem) Next() *ListItem {
	if item == nil {
		return nil
	}
	return item.next
}

func (item *ListItem) remove() {
	if item == nil {
		return
	}
	if item.prev != nil {
		item.prev.next = item.next
	}
	if item.next != nil {
		item.next.prev = item.prev
	}
}

// ListStoreImpl implement the ListStore interface.
// We will store all the objects in hashmap and linked-list, and maintain the StoredObj's generation
// according to ListStoreImpl.generation.
// And the linked-list will be organized in descending order of generation.
type ListStoreImpl struct {
	store      map[string]*ListItem
	head       *ListItem
	generation uint64 // generation hold the global max-generation in this store.
}

var (
	_ Store     = &ListStoreImpl{}
	_ ListStore = &ListStoreImpl{}
)

func NewListStore() ListStore {
	return &ListStoreImpl{
		store: make(map[string]*ListItem),
	}
}

func (s *ListStoreImpl) Get(key string) StoredObj {
	if s == nil {
		return nil
	}
	item, ok := s.store[key]
	if !ok || item == nil {
		return nil
	}
	return item.StoredObj
}

// Set will update the ListItem if it has been existed, or create a new ListItem to hold it.
// The ListStoreImpl.generation will be updated and the ListItem will be moved to head.
func (s *ListStoreImpl) Set(key string, obj StoredObj) {
	if s == nil || obj == nil {
		// obj == nil, return directly.
		return
	}
	var item *ListItem
	if item = s.store[key]; item != nil {
		s.remove(item)
		item.StoredObj = obj
	} else {
		item = newListItem(key, obj)
	}
	obj.SetGeneration(nextGeneration(&s.generation))
	if s.head != nil {
		s.head.prev = item
	}
	item.next, item.prev = s.head, nil
	s.head = item
	s.store[key] = item
}

func (s *ListStoreImpl) Delete(key string) {
	if s == nil {
		return
	}
	if item, ok := s.store[key]; ok && item != nil {
		s.remove(item)
	}
	delete(s.store, key)
}

func (s *ListStoreImpl) Front() *ListItem {
	if s == nil {
		return nil
	}
	return s.head
}

func (s *ListStoreImpl) Len() int {
	if s == nil {
		return 0
	}
	return len(s.store)
}

// UpdateRawStore update RawStore according the generation.
// We will update the RawStore by linked-list firstly, and by RawStore.UpdatedSet. Finally refresh the
// RawStore.generation and do cleanup.
func (s *ListStoreImpl) UpdateRawStore(store RawStore, cloneFunc CloneFunc, cleanFunc CleanFunc) {
	if s == nil || store == nil {
		return
	}
	storedGeneration := store.GetGeneration()
	for e := s.Front(); e != nil; e = e.Next() {
		if e.GetGeneration() <= storedGeneration {
			break
		}
		cloneFunc(e.key, e.StoredObj)
	}
	for _, key := range store.UpdatedSet().UnsortedList() {
		if s.store[key] != nil {
			cloneFunc(key, s.store[key].StoredObj)
		}
	}
	store.SetGeneration(s.generation)
	store.ResetUpdatedSet()
	cleanFunc()
}

func (s *ListStoreImpl) HashStore() HashStore {
	if s == nil {
		return nil
	}
	ret := make(HashStore, s.Len())
	for k, item := range s.store {
		ret[k] = item.StoredObj
	}
	return ret
}

func (s *ListStoreImpl) String() string {
	if s == nil {
		return "{}"
	}
	items := []string{}
	for k, item := range s.store {
		items = append(items, fmt.Sprintf("{%v:%v}", k, item.StoredObj))
	}
	sort.Strings(items)
	return fmt.Sprintf("{Store:%v}", items)
}

func (s *ListStoreImpl) remove(item *ListItem) {
	if s == nil || item == nil {
		return
	}
	item.remove()
	if s.head == item {
		s.head = item.Next()
	}
}
