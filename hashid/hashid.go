package hashid

import (
	"sync"
)

type HashID interface {
	Get(any) int64
	Lookup(int64) any
	Exist(any) bool
	Len() int64
}

type HashIDImpl struct {
	data   map[any]int64
	lookup map[int64]any
	id     int64
	mu     sync.RWMutex
}

func NewHashID() HashID {
	return &HashIDImpl{
		data:   make(map[any]int64),
		lookup: make(map[int64]any),
		id:     0,
	}
}

func (h *HashIDImpl) Get(key any) int64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.data[key]
	if !ok {
		h.data[key] = h.id
		h.lookup[h.id] = key
		h.id++
	}
	return h.data[key]
}

func (h *HashIDImpl) Lookup(id int64) any {
	h.mu.RLock()
	defer h.mu.RUnlock()
	key := h.lookup[id]
	return key
}

func (h *HashIDImpl) Exist(key any) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.data[key]
	return ok
}

func (h *HashIDImpl) Len() int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.id
}
