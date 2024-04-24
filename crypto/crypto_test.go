package crypto

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {

	payload := "some file that needs to be encrypted"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := NewEncryptionKey()

	_, err := CopyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Encrypted: %s\n", dst.String())

	output := new(bytes.Buffer)
	if _, err := CopyDecrypt(key, dst, output); err != nil {
		t.Error(err)
	}

	fmt.Printf("Decrypted: %s\n", output.String())

	if output.String() != payload {
		t.Errorf("Decryption failed")
	}

}
