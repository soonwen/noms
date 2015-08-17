package types

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/attic-labs/noms/chunks"
	"github.com/attic-labs/noms/enc"
	"github.com/attic-labs/noms/ref"
	"github.com/stretchr/testify/assert"
)

type ErrorSource struct {
	chunks.ChunkSource
}

func (e ErrorSource) Get(r ref.Ref) (io.ReadCloser, error) {
	return nil, errors.New("Good golly Miss Molly!")
}

func TestTolerateUngettableRefs(t *testing.T) {
	assert := assert.New(t)
	v, err := ReadValue(ref.Ref{}, &chunks.TestStore{})
	assert.Nil(v)
	assert.Nil(err)

	v, err = ReadValue(ref.Ref{}, &ErrorSource{})
	assert.Nil(v)
	assert.NotNil(err)
}

func TestBlobLeafDecode(t *testing.T) {
	assert := assert.New(t)

	blobLeafDecode := func(r io.Reader) Value {
		i, err := enc.Decode(r)
		assert.NoError(err)
		f, err := fromEncodeable(i, nil)
		assert.NoError(err)
		val, err := f.Deref(nil)
		assert.NoError(err)
		return val
	}

	reader := bytes.NewBufferString("b ")
	v1 := blobLeafDecode(reader)
	bl1 := newBlobLeaf([]byte{})
	assert.True(bl1.Equals(v1))

	reader = bytes.NewBufferString("b Hello World!")
	v2 := blobLeafDecode(reader)
	bl2 := newBlobLeaf([]byte("Hello World!"))
	assert.True(bl2.Equals(v2))
}