package des

import (
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// 3DES-CBC 加密数据
func CBCTripleEncrypt(originData, key []byte) ([]byte, error) {
	return cbcTripleEncrypt(originData, key)
}

// 3DES-CBC 解密数据
func CBCTripleDecrypt(secretData, key []byte) ([]byte, error) {
	return cbcTripleDecrypt(secretData, key)
}

// DES-CBC 加密数据
func CBCEncrypt(originData, key []byte) ([]byte, error) {
	return cbcEncrypt(originData, key)
}

// DES-CBC 解密数据
func CBCDecrypt(secretData, key []byte) ([]byte, error) {
	return cbcDecrypt(secretData, key)
}

func cbcTripleEncrypt(originData, key []byte) ([]byte, error) {
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

func cbcTripleDecrypt(secretData, key []byte) (originByte []byte, err error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originByte = make([]byte, len(secretData))
	blockMode.CryptBlocks(originByte, secretData)
	if len(originByte) == 0 {
		return nil, errors.New("blockMode.CryptBlocks error")
	}
	return PKCS7UnPadding(originByte), nil
}

func cbcEncrypt(originData, key []byte) ([]byte, error) {
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

func cbcDecrypt(secretData, key []byte) (originByte []byte, err error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originByte = make([]byte, len(secretData))
	blockMode.CryptBlocks(originByte, secretData)
	if len(originByte) == 0 {
		return nil, errors.New("blockMode.CryptBlocks error")
	}
	return PKCS7UnPadding(originByte), nil
}
