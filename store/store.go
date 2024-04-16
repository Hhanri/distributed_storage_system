package store

import (
	"bytes"
	"io"
	"log"
	"os"
)

const defaultRootStorage = "../storage"

type StoreOpts struct {
	// root folder for storage
	root          string
	pathTransform pathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.pathTransform == nil {
		opts.pathTransform = defaultPathTransform
	}
	if opts.root == "" {
		opts.root = defaultRootStorage
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) has(key string) bool {
	pathKey := s.pathTransform(s.root, key)

	file, err := os.Stat(pathKey.FullPath())
	if err != nil {
		return false
	}

	return file != nil
}

func (s *Store) delete(key string) error {
	pathKey := s.pathTransform(s.root, key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	return os.RemoveAll(pathKey.FullPath())
}

func (s *Store) read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buff := new(bytes.Buffer)

	_, err = io.Copy(buff, f)

	return buff, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.pathTransform(s.root, key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) writeStream(key string, reader io.Reader) error {

	pathKey := s.pathTransform(s.root, key)

	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	fullPath := pathKey.FullPath()

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(file, reader)
	if err != nil {
		return err
	}

	log.Printf("Written (%d) bytes to disk: %s", n, fullPath)

	return nil
}
