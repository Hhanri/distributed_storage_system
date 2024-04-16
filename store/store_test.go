package store

import (
	"bytes"
	"io"
	"testing"
)

func newTestStore() *Store {
	opts := StoreOpts{
		pathTransform: hashPathTransform,
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

	key := "myImageKey"

	data := []byte("some jpg bytes idk just go with it")

	if err := store.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if !store.has(key) {
		t.Errorf("no file found for key: %s", key)
	}

	reader, err := store.read(key)
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
}

func TestDelete(t *testing.T) {

	store := newTestStore()

	defer tearDown(t, store)

	key := "myImageKey"

	data := []byte("some jpg bytes idk just go with it")

	if err := store.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := store.delete(key); err != nil {
		t.Error(err)
	}

}
