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
	// Root folder for storage
	Root          string
	PathTransform pathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransform == nil {
		opts.PathTransform = defaultPathTransform
	}
	if opts.Root == "" {
		opts.Root = DefaultRootStorage
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransform(key)

	_, err := os.Stat(pathKey.FullPath(s.Root))
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransform(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	return os.RemoveAll(pathKey.FullPath(s.Root))
}

func (s *Store) clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Read(key string) (int64, io.Reader, error) {
	n, f, err := s.readStream(key)
	if err != nil {
		return n, nil, err
	}

	defer f.Close()

	buff := new(bytes.Buffer)

	_, err = io.Copy(buff, f)

	return n, buff, err
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransform(key)
	file, err := os.Open(pathKey.FullPath(s.Root))
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}
	return fi.Size(), file, nil
}

func (s *Store) Write(key string, reader io.Reader) (int64, error) {
	return s.writeStream(key, reader)
}

func (s *Store) writeStream(key string, reader io.Reader) (int64, error) {

	pathKey := s.PathTransform(key)

	if err := os.MkdirAll(pathKey.DirPath(s.Root), os.ModePerm); err != nil {
		return 0, err
	}

	fullPath := pathKey.FullPath(s.Root)

	file, err := os.Create(fullPath)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(file, reader)
	if err != nil {
		return 0, err
	}

	log.Printf("Written (%d) bytes to disk: %s", n, fullPath)

	return n, nil
}
