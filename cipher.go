package tapo

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type Cipher struct {
	key []byte
	iv  []byte
}

func Pad(buf []byte, size int) ([]byte, error) {
	bufLen := len(buf)
	padLen := size - bufLen%size
	padded := make([]byte, bufLen+padLen)
	copy(padded, buf)
	for i := 0; i < padLen; i++ {
		padded[bufLen+i] = byte(0)
	}
	return padded, nil
}

func Unpad(padded []byte, size int) ([]byte, error) {
	if len(padded)%size != 0 {
		return nil, errors.New("pkcs7: Padded value wasn't in correct size.")
	}

	bufLen := len(padded) - int(padded[len(padded)-1])
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf, nil
}

func (c *Cipher) Encrypt(payload []byte) []byte {
	block, _ := aes.NewCipher(c.key)
	encrypter := cipher.NewCBCEncrypter(block, c.iv)

	paddedPayload, _ := Pad(payload, aes.BlockSize)
	encryptedPayload := make([]byte, len(paddedPayload))
	encrypter.CryptBlocks(encryptedPayload, paddedPayload)

	return encryptedPayload
}

func (c *Cipher) Decrypt(payload []byte) []byte {
	block, _ := aes.NewCipher(c.key)
	encrypter := cipher.NewCBCDecrypter(block, c.iv)

	decryptedPayload := make([]byte, len(payload))
	encrypter.CryptBlocks(decryptedPayload, payload)

	unpaddedPayload, _ := Unpad(decryptedPayload, aes.BlockSize)

	return unpaddedPayload
}
