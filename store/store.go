package store

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"
)

const DefaultRootStorage string = "../storage"

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
		opts.root = DefaultRootStorage
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.pathTransform(key)

	_, err := os.Stat(pathKey.FullPath(s.root))
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathKey := s.pathTransform(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	return os.RemoveAll(pathKey.FullPath(s.root))
}

func (s *Store) clear() error {
	return os.RemoveAll(s.root)
}

func (s *Store) Read(key string) (io.Reader, error) {
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
	pathKey := s.pathTransform(key)
	return os.Open(pathKey.FullPath(s.root))
}

func (s *Store) Write(key string, reader io.Reader) error {
	return s.writeStream(key, reader)
}

func (s *Store) writeStream(key string, reader io.Reader) error {

	pathKey := s.pathTransform(key)

	if err := os.MkdirAll(pathKey.DirPath(s.root), os.ModePerm); err != nil {
		return err
	}

	fullPath := pathKey.FullPath(s.root)

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
