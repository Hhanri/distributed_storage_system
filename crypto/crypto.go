package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func HashKey(key string) string {
	hash := md5.Sum([]byte(key))

	return hex.EncodeToString(hash[:])
}

func NewEncryptionKey() []byte {
	keyBuff := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuff)
	return keyBuff
}

func CopyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend the IV to the file
	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)

	return copyStream(stream, block.BlockSize(), src, dst)
}

func CopyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// read the IV from the given src
	// should be in our case block.BlockSIze() bytes to read
	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)

	return copyStream(stream, block.BlockSize(), src, dst)
}

func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {
	buff := make([]byte, 32*1024)
	nw := blockSize

	for {
		n, err := src.Read(buff)

		if n > 0 {
			stream.XORKeyStream(buff, buff[:n])

			nn, err := dst.Write(buff[:n])
			if err != nil {
				return 0, err
			}

			nw += nn
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}

	return nw, nil
}
