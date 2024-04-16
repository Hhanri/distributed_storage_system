package store

import (
	"bytes"
	"io"
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

	opts := StoreOpts{
		pathTransform: hashPathTransform,
	}
	store := NewStore(opts)

	key := "myImageKey"

	data := []byte("some jpg bytes idk just go with it")

	if err := store.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := store.delete(key); err != nil {
		t.Error(err)
	}

}
