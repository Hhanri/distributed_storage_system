package store

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newTestStore() *Store {
	opts := StoreOpts{
		PathTransform: hashPathTransform,
	}
	store := NewStore(opts)
	return store
}

func tearDown(t *testing.T, store *Store) {
	if err := store.clear(); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {

	store := newTestStore()

	defer tearDown(t, store)

	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("myImageKey_%d", i)

		data := []byte("some jpg bytes idk just go with it")

		if err := store.Write(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if !store.Has(key) {
			t.Errorf("no file found for key: %s", key)
		}

		reader, err := store.Read(key)
		if err != nil {
			t.Error(err)
		}

		bytes, err := io.ReadAll(reader)
		if err != nil {
			t.Error(err)
		}

		if string(bytes) != string(data) {
			t.Errorf("expected %s\ngot %s", data, bytes)
		}

		if err := store.Delete(key); err != nil {
			t.Error(err)
		}

		if store.Has(key) {
			t.Errorf("should not have key %s", key)
		}

	}

}
