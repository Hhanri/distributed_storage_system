package store

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {

	opts := StoreOpts{
		pathTransform: hashPathTransform,
	}
	store := NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes idk just go with it"))

	if err := store.writeStream("myImageKey", data); err != nil {
		t.Error(err)
	}

}
