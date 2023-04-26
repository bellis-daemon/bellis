package cryptoo

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type ConnAes struct {
	key      string
	keyBlock cipher.Block
	iv       string
}

func NewConnAes(key string, iv string) (ConnAes, error) {
	if len(iv) != aes.BlockSize {
		return ConnAes{}, errors.New("errors iv size not 16")
	}
	keyBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ConnAes{}, err
	}
	return ConnAes{
		key:      key,
		keyBlock: keyBlock,
		iv:       iv,
	}, nil
}

func (ca *ConnAes) Encrypt(src []byte) []byte {
	paddingLen := aes.BlockSize - (len(src) % aes.BlockSize)
	for i := 0; i < paddingLen; i++ {
		src = append(src, byte(paddingLen))
	}
	enbuf := make([]byte, len(src))
	cbce := cipher.NewCBCEncrypter(ca.keyBlock, []byte(ca.iv))
	cbce.CryptBlocks(enbuf, src)
	return enbuf
}

func (ca *ConnAes) Decrypt(src []byte) ([]byte, error) {
	if (len(src) < aes.BlockSize) || (len(src)%aes.BlockSize != 0) {
		return nil, errors.New("errors encrypt data size1")
	}
	debuf := make([]byte, len(src))
	cbcd := cipher.NewCBCDecrypter(ca.keyBlock, []byte(ca.iv))
	cbcd.CryptBlocks(debuf, src)
	padding := int(debuf[len(src)-1])
	// glg.Debug(debuf)
	if padding > 16 {
		return nil, errors.New("errors encrypt data size2")
	}
	return debuf[:len(src)-padding], nil
}
