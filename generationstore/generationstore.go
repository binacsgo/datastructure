package generationstore

import (
	"sync/atomic"

	"k8s.io/apimachinery/pkg/util/sets"
)

func nextGeneration(generation *int64) int64 {
	return atomic.AddInt64(generation, 1)
}

// ------------------- StoredObj -------------------
type StoredObj interface {
	GetGeneration() int64
	SetGeneration(int64)
}

// ------------------- Store -------------------
type Store interface {
	Get(string) StoredObj
	Set(string, StoredObj)
	Delete(string)
	Len() int
	HashStore() HashStore
}

// ------------------- ListStore -------------------
type ListStore interface {
	Store
	Front() *ListItem
	UpdateRawStore(RawStore, CloneFunc, CleanFunc)
}

// ------------------- RawStore -------------------
type RawStore interface {
	Store
	SetGeneration(int64)
	GetGeneration() int64
	UpdatedSet() sets.String
}

type (
	CloneFunc func(string, StoredObj)
	CleanFunc func()

	HashStore map[string]StoredObj
)
