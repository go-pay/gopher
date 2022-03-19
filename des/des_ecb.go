package des

import (
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// 3DES-ECB 加密数据
func ECBTripleEncrypt(originData, key, iv []byte) ([]byte, error) {
	data, err := ecbTripleEncrypt(originData, key, iv)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 3DES-ECB 解密数据
func ECBTripleDecrypt(secretData, key, iv []byte) ([]byte, error) {
	return ecbTripleDecrypt(secretData, key, iv)
}

// DES-ECB 加密数据
func ECBEncrypt(originData, key, iv []byte) ([]byte, error) {
	data, err := ecbEncrypt(originData, key, iv)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DES-ECB 解密数据
func ECBDecrypt(secretData, key, iv []byte) ([]byte, error) {
	return ecbDecrypt(secretData, key, iv)
}

func ecbTripleEncrypt(originData, key, iv []byte) ([]byte, error) {
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

func ecbTripleDecrypt(secretData, desKey, iv []byte) (originByte []byte, err error) {
	block, err := des.NewTripleDESCipher(desKey)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
	originByte = make([]byte, len(secretData))
	blockMode.CryptBlocks(originByte, secretData)
	if len(originByte) == 0 {
		return nil, errors.New("blockMode.CryptBlocks error")
	}
	return PKCS7UnPadding(originByte), nil
}

func ecbEncrypt(originData, key, iv []byte) ([]byte, error) {
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

func ecbDecrypt(secretData, desKey, iv []byte) (originByte []byte, err error) {
	block, err := des.NewCipher(desKey)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:block.BlockSize()])
	originByte = make([]byte, len(secretData))
	blockMode.CryptBlocks(originByte, secretData)
	if len(originByte) == 0 {
		return nil, errors.New("blockMode.CryptBlocks error")
	}
	return PKCS7UnPadding(originByte), nil
}
