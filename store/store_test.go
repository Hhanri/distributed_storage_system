package store

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestStore(t *testing.T) {

	opts := StoreOpts{
		pathTransform: hashPathTransform,
	}
	store := NewStore(opts)

	key := "myImageKey"

	data := []byte("some jpg bytes idk just go with it")

	if err := store.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	reader, err := store.read(key)
	if err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != string(data) {
		t.Errorf("expected %s\ngot %s", data, bytes)
	}
}
