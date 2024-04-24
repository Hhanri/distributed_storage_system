package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

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

	buff := make([]byte, 32*1024)
	stream := cipher.NewCTR(block, iv)

	for {
		n, err := src.Read(buff)
		if n > 0 {
			stream.XORKeyStream(buff, buff[:n])
			if _, err := dst.Write(buff[:n]); err != nil {
				return 0, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}
	return 0, nil
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

	buff := make([]byte, 32*1024)
	stream := cipher.NewCTR(block, iv)

	for {
		n, err := src.Read(buff)

		if n > 0 {
			stream.XORKeyStream(buff, buff[:n])
			if _, err := dst.Write(buff[:n]); err != nil {
				return 0, err
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
