package generationstore

import "k8s.io/apimachinery/pkg/util/sets"

// StoredObj defines the methods that all the objects stored in the generationstore must have.
type StoredObj interface {
	GetGeneration() uint64
	SetGeneration(uint64)
}

type (
	StoreRangeFunc          func(string, StoredObj)
	StoreConditionRangeFunc func(string, StoredObj) bool
)

// Store defines the field that the some-datastructure will hold if it need generatestore.
// The Store field's real object is either ListStore or RawStore.
type Store interface {
	Get(string) StoredObj
	Set(string, StoredObj)
	Delete(string)
	Len() int
	Range(StoreRangeFunc)
	ConditionRange(StoreConditionRangeFunc) bool
	String() string
}

// ListStore defines the methods of ListStore.
type ListStore interface {
	Store
	Front() *ListItem
	UpdateRawStore(RawStore, CloneFunc, CleanFunc)
}

// RawStore defines the methods of RawStore.
type RawStore interface {
	Store
	SetGeneration(uint64)
	GetGeneration() uint64
	UpdatedSet() sets.String
	ResetUpdatedSet()
}

type (
	// CloneFunc defines the clone function used in UpdateRawStore.
	CloneFunc func(string, StoredObj)
	// CleanFunc defines the cleanup function used in UpdateRawStore.
	CleanFunc func()

	// HashStore defines the data model used for traversal.
	HashStore map[string]StoredObj
)

func DefaultCleanFunc(cache ListStore, snapshot RawStore) CleanFunc {
	return func() {
		if cache == nil || snapshot == nil {
			return
		}
		if cache.Len() != snapshot.Len() {
			diff := snapshot.Len() - cache.Len()
			snapshot.ConditionRange(
				func(key string, _ StoredObj) bool {
					if diff <= 0 {
						return false
					}
					if cache.Get(key) == nil {
						snapshot.Delete(key)
						diff--
					}
					return true
				})
		}
	}
}
