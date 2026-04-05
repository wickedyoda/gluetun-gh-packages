package pkcs8

import (
	"crypto/cipher"
	"crypto/des" //nolint:gosec
	"encoding/asn1"
	"errors"
	"fmt"

	pkcs8lib "github.com/youmark/pkcs8"
)

var ErrEncryptDESCBCUnsupported = errors.New("encrypting private keys with DES-CBC is unsupported")

func init() { //nolint:gochecknoinits
	pkcs8lib.RegisterCipher(oidDESCBC, newCipherDESCBCBlock)
}

func newCipherDESCBCBlock() pkcs8lib.Cipher {
	return cipherDESCBC{}
}

type cipherDESCBC struct{}

func (c cipherDESCBC) IVSize() int {
	return des.BlockSize
}

func (c cipherDESCBC) KeySize() int {
	return 8 //nolint:mnd
}

func (c cipherDESCBC) OID() asn1.ObjectIdentifier {
	return oidDESCBC
}

func (c cipherDESCBC) Encrypt(_, _, _ []byte) ([]byte, error) {
	return nil, ErrEncryptDESCBCUnsupported
}

func (c cipherDESCBC) Decrypt(key, iv, ciphertext []byte) ([]byte, error) {
	block, err := des.NewCipher(key) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("creating DES cipher: %w", err)
	}
	blockDecrypter := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	blockDecrypter.CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}
