package store

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

type pathTransformFunc func(root string, key string) PathKey

var defaultPathTransform = func(root string, key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

func hashPathTransform(root string, key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLength := len(hashStr) / blockSize

	paths := make([]string, sliceLength)

	for i := 0; i < sliceLength; i++ {
		from, to := i*blockSize, i*blockSize+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		FileName: hashStr,
		PathName: root + "/" + strings.Join(paths, "/"),
	}
}

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}
