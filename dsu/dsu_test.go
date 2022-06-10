package dsu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSUImpl(t *testing.T) {
	n := int64(10)
	dsu := NewDSU(n)
	for i := int64(0); i <= n; i++ {
		assert.Equal(t, i, dsu.Find(i))
	}

	dsu.Merge(2, 3)
	assert.Equal(t, dsu.Find(2), dsu.Find(3))
	dsu.Merge(4, 5)
	assert.Equal(t, dsu.Find(4), dsu.Find(5))

	assert.NotEqual(t, dsu.Find(2), dsu.Find(5))
	dsu.Merge(2, 4)
	assert.Equal(t, dsu.Find(2), dsu.Find(5))

	assert.Equal(t, int64(4), dsu.SetSize(2))
	assert.Equal(t, int64(7), dsu.SCC())

	dsu.Squash()
}
