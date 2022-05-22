package generationstore

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

// ------------------- ListStoreImpl -------------------
type ListStoreImpl struct {
	store      map[string]*ListItem
	head       *ListItem
	generation int64
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

func (s *ListStoreImpl) Set(key string, obj StoredObj) {
	if s == nil {
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

func (s *ListStoreImpl) remove(item *ListItem) {
	if s == nil || item == nil {
		return
	}
	item.remove()
	if s.head == item {
		s.head = item.Next()
	}
}

func DefaultCleanFunc(cache ListStore, snapshot RawStore) CleanFunc {
	return func() {
		if cache == nil || snapshot == nil {
			return
		}
		if cache.Len() != snapshot.Len() {
			diff := snapshot.Len() - cache.Len()
			for key := range snapshot.HashStore() {
				if diff <= 0 {
					break
				}
				if cache.Get(key) == nil {
					snapshot.Delete(key)
					diff--
				}
			}
		}
	}
}
