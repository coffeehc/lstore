package mheap

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func Test_allocate(t *testing.T) {
	should := require.New(t)
	mgr, err := New(8, 2)
	should.NoError(err)
	allocated := mgr.Allocate(7, []byte{1, 2, 3})
	should.Equal([]byte{1, 2, 3}, allocated)
}