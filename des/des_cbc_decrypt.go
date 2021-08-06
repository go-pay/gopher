package des

import (
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// 解密数据的Bytes数组
func CBCTripleDecryptData(secretData, key []byte) ([]byte, error) {
	return tripleDecrypt(secretData, key)
}

// 解密数据的Bytes数组
func CBCDecryptData(secretData, key []byte) ([]byte, error) {
	return decrypt(secretData, key)
}

// 解密数据的Bytes数组
func CBCTripleDecryptIvData(secretData, key, iv []byte) ([]byte, error) {
	return tripleDecryptIv(secretData, key, iv)
}

// 解密数据的Bytes数组
func CBCDecryptIvData(secretData, key, iv []byte) ([]byte, error) {
	return decryptIv(secretData, key, iv)
}

func tripleDecrypt(secretData, key []byte) (originByte []byte, err error) {
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

func decrypt(secretData, key []byte) (originByte []byte, err error) {
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

func tripleDecryptIv(secretData, desKey, iv []byte) (originByte []byte, err error) {
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

func decryptIv(secretData, desKey, iv []byte) (originByte []byte, err error) {
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
