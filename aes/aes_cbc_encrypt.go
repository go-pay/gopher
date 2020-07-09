package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// 加密后转成Base64字符串
func AesCBCEncryptToString(jsonData []byte, secretKey string) (string, error) {
	key := []byte(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	jsonData = PKCS5Padding(jsonData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(jsonData))
	blockMode.CryptBlocks(crypted, jsonData)
	secretData := base64.StdEncoding.EncodeToString(crypted)
	return secretData, nil
}

// 加密后的Bytes数组
func AesCBCEncryptToBytes(jsonData []byte, secretKey string) ([]byte, error) {
	key := []byte(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	jsonData = PKCS5Padding(jsonData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	secretData := make([]byte, len(jsonData))
	blockMode.CryptBlocks(secretData, jsonData)
	return secretData, nil
}
