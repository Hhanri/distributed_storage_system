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
	return s.readStream(key)
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

func (s *Store) WriteDecrypt(encKey []byte, key string, reader io.Reader) (int64, error) {
	file, err := s.openFileForWriting(key, reader)
	if err != nil {
		return 0, err
	}

	n, err := crypto.CopyDecrypt(encKey, reader, file)
	return int64(n), err
}

func (s *Store) writeStream(key string, reader io.Reader) (int64, error) {
	file, err := s.openFileForWriting(key, reader)
	if err != nil {
		return 0, err
	}

	return io.Copy(file, reader)

}

func (s *Store) LogWrite(n int64, addr string) {
	log.Printf("Written (%d) bytes to: %s", n, addr)
}

func (s *Store) openFileForWriting(key string, reader io.Reader) (*os.File, error) {
	pathKey := s.PathTransform(key)

	if err := os.MkdirAll(pathKey.DirPath(s.Root), os.ModePerm); err != nil {
		return nil, err
	}

	fullPath := pathKey.FullPath(s.Root)

	return os.Create(fullPath)
}
