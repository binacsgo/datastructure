package generationstore

import (
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/util/sets"
)

// RawStoreImpl implement the RawStore interface.
// StoredObj will be stored in a raw hashmap, and RawStoreImpl.generation hold
// the max-generation of these items.
type RawStoreImpl struct {
	store      HashStore
	generation uint64
	updatedSet sets.String // updatedSet record all the items may be changed in RawStoreImpl.
}

var (
	_ Store    = &RawStoreImpl{}
	_ RawStore = &RawStoreImpl{}
)

func NewRawStore() RawStore {
	return &RawStoreImpl{
		store:      make(HashStore),
		updatedSet: sets.NewString(),
	}
}

func (s *RawStoreImpl) Get(key string) StoredObj {
	if s == nil {
		return nil
	}
	item, ok := s.store[key]
	if !ok || item == nil {
		return nil
	}
	return item
}

func (s *RawStoreImpl) Set(key string, obj StoredObj) {
	if s == nil {
		return
	}
	s.store[key] = obj
	s.updatedSet.Insert(key)
}

func (s *RawStoreImpl) Delete(key string) {
	if s == nil {
		return
	}
	delete(s.store, key)
	s.updatedSet.Insert(key)
}

func (s *RawStoreImpl) Len() int {
	if s == nil {
		return 0
	}
	return len(s.store)
}

func (s *RawStoreImpl) HashStore() HashStore {
	if s == nil {
		return nil
	}
	ret := make(HashStore, s.Len())
	for k := range s.store {
		ret[k] = s.store[k]
	}
	return ret
}

func (s *RawStoreImpl) SetGeneration(generation uint64) {
	if s == nil {
		return
	}
	s.generation = generation
}

func (s *RawStoreImpl) GetGeneration() uint64 {
	if s == nil {
		return 0
	}
	return s.generation
}

func (s *RawStoreImpl) UpdatedSet() sets.String {
	if s == nil {
		return sets.NewString()
	}
	return s.updatedSet
}

func (s *RawStoreImpl) ResetUpdatedSet() {
	if s == nil {
		return
	}
	s.updatedSet = sets.NewString()
}

func (s *RawStoreImpl) String() string {
	if s == nil {
		return "{}"
	}
	items := []string{}
	for k, item := range s.store {
		items = append(items, fmt.Sprintf("{%v:%v}", k, item))
	}
	sort.Strings(items)
	return fmt.Sprintf("{Store:%v}", items)
}
