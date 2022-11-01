package splay

type Comparable interface {
	// Compare defines a comparison function in splay that returns `true` if and only if the
	// current element is strictly greater than the incoming element.
	Compare(Comparable) bool
}

type MaintainInfo interface {
	// Maintain defines the maintenance operation in the splay, which contains the properties
	// of the subtree rooted at the current node. We will update the properties of the current
	// node based on its left and right children.
	Maintain(MaintainInfo, MaintainInfo)

	// Clone return a clone of the MaintainInfo
	Clone() MaintainInfo

	// String implements the String interface.
	String() string
}

// StoredObj defines all the methods that need to be implemented by the element being stored.
type StoredObj interface {
	// Key returns the unique key used by the object in the splay.
	Key() string
	// String implements the String interface.
	String() string
	// MakeMaintainInfo return a MaintainInfo used by the StoredObj.
	MakeMaintainInfo() MaintainInfo

	Comparable
}

type (
	// RangeFunc visit objects by inorder traversal.
	RangeFunc func(StoredObj)
	// ConditionRangeFunc visit objects by inorder traversal.
	ConditionRangeFunc func(StoredObj) bool
)

// Splay defines all methods of the splay-tree.
type Splay interface {
	// Insert a StoredObj into the splay. Returns true if successful.
	Insert(StoredObj) bool
	// Delete a StoredObj from the splay. Returns true if successful.
	Delete(StoredObj) bool
	// Get a StoredObj from the splay.
	Get(StoredObj) StoredObj
	// Partition will bring together all objects strictly smaller than the current object
	// in a subtree and return the root of the subtree.
	Partition(Comparable) StoredObj
	// Range traverses the entire splay in mid-order.
	Range(RangeFunc)
	// ConditionRange traverses the entire splay in mid-order and ends the access immediately
	// if ConditionRangeFunc returns false.
	ConditionRange(ConditionRangeFunc)
	// Len returns the number of all objects in the splay.
	Len() int
	// String implements the String interface.
	String() string
	// Clone return a clone of the Splay.
	Clone() Splay
	// PrintTree outputs splay in the form of a tree diagram.
	PrintTree() string
}

// maintainInfoForLookup defines one of the simplest MaintainInfo implementations for lookups only.
type maintainInfoForLookup struct{}

func (o *maintainInfoForLookup) Maintain(l, r MaintainInfo) {}
func (o *maintainInfoForLookup) Clone() MaintainInfo        { return &maintainInfoForLookup{} }
func (o *maintainInfoForLookup) String() string             { return "maintainInfoForLookup" }

// storedObjForLookup defines one of the simplest StoredObj implementations for lookups only.
type storedObjForLookup struct{ key string }

func (o *storedObjForLookup) Key() string                    { return o.key }
func (o *storedObjForLookup) String() string                 { return o.key }
func (o *storedObjForLookup) Maintain(_, _ StoredObj)        {}
func (o *storedObjForLookup) MakeMaintainInfo() MaintainInfo { return &maintainInfoForLookup{} }
func (o *storedObjForLookup) Compare(Comparable) bool        { return false }
func NewStoredObjForLookup(key string) StoredObj {
	return &storedObjForLookup{
		key: key,
	}
}

var (
	_ MaintainInfo = &maintainInfoForLookup{}
	_ StoredObj    = &storedObjForLookup{}

	NilObj = NewStoredObjForLookup("NilObj")
	MinObj = NewStoredObjForLookup("MinObj")
	MaxObj = NewStoredObjForLookup("MaxObj")
)
