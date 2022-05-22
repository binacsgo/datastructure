package generationstore

import "k8s.io/apimachinery/pkg/util/sets"

// ------------------- RawStoreImpl -------------------
type RawStoreImpl struct {
	store      HashStore
	generation int64
	updatedSet sets.String
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
	return s.store
}

func (s *RawStoreImpl) SetGeneration(generation int64) {
	if s == nil {
		return
	}
	s.generation = generation
}

func (s *RawStoreImpl) GetGeneration() int64 {
	if s == nil {
		return 0
	}
	return s.generation
}

func (s *RawStoreImpl) UpdatedSet() sets.String {
	return s.updatedSet
}
