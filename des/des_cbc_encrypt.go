package des

import (
	"crypto/cipher"
	"crypto/des"
)

// 加密后的Bytes数组
func CBCTripleEncryptData(originData, key []byte) ([]byte, error) {
	return tripleEncrypt(originData, key)
}

// 加密后的Bytes数组
func CBCEncryptData(originData, key []byte) ([]byte, error) {
	return encrypt(originData, key)
}

// 加密后的Bytes数组
func CBCTripleEncryptIvData(originData, key, iv []byte) ([]byte, error) {
	data, err := tripleEncryptIv(originData, key, iv)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 加密后的Bytes数组
func CBCEncryptIvData(originData, key, iv []byte) ([]byte, error) {
	data, err := encryptIv(originData, key, iv)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func tripleEncrypt(originData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	originData = PKCS7Padding(originData, blockSize)
	secretData := make([]byte, len(originData))
	blockMode.CryptBlocks(secretData, originData)
	return secretData, nil
}

func encrypt(originData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	originData = PKCS7Padding(originData, blockSize)
	secretData := make([]byte, len(originData))
	blockMode.CryptBlocks(secretData, originData)
	return secretData, nil
}

func tripleEncryptIv(originData, key, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])

	originData = PKCS7Padding(originData, block.BlockSize())
	secretData := make([]byte, len(originData))
	blockMode.CryptBlocks(secretData, originData)
	return secretData, nil
}

func encryptIv(originData, key, iv []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCEncrypter(block, iv[:block.BlockSize()])

	originData = PKCS7Padding(originData, block.BlockSize())
	secretData := make([]byte, len(originData))
	blockMode.CryptBlocks(secretData, originData)
	return secretData, nil
}
