package store

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/Hhanri/distributed_storage_system/crypto"
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

func (s *Store) Has(id, key string) bool {
	pathKey := s.PathTransform(key)

	_, err := os.Stat(pathKey.FullPath(s.Root, id))
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Delete(id, key string) error {
	pathKey := s.PathTransform(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	return os.RemoveAll(pathKey.FullPath(s.Root, id))
}

func (s *Store) clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Read(id, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

func (s *Store) readStream(id, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransform(key)
	file, err := os.Open(pathKey.FullPath(s.Root, id))
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}
	return fi.Size(), file, nil
}

func (s *Store) Write(id, key string, reader io.Reader) (int64, error) {
	return s.writeStream(id, key, reader)
}

func (s *Store) WriteDecrypt(id string, encKey []byte, key string, reader io.Reader) (int64, error) {
	file, err := s.openFileForWriting(id, key, reader)
	if err != nil {
		return 0, err
	}

	n, err := crypto.CopyDecrypt(encKey, reader, file)
	return int64(n), err
}

func (s *Store) writeStream(id, key string, reader io.Reader) (int64, error) {
	file, err := s.openFileForWriting(id, key, reader)
	if err != nil {
		return 0, err
	}

	return io.Copy(file, reader)

}

func (s *Store) LogWrite(n int64, addr string) {
	log.Printf("Written (%d) bytes to: %s", n, addr)
}

func (s *Store) openFileForWriting(id, key string, reader io.Reader) (*os.File, error) {
	pathKey := s.PathTransform(key)

	if err := os.MkdirAll(pathKey.DirPath(s.Root, id), os.ModePerm); err != nil {
		return nil, err
	}

	fullPath := pathKey.FullPath(s.Root, id)

	return os.Create(fullPath)
}
