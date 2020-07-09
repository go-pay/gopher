package xrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// RSA加密数据
// originData：原始字符串
// publicKeyFilePath：公钥证书文件路径
func RsaEncryptData(originData string, publicKeyFilePath string) (cipherData string, err error) {
	fileBytes, err := ioutil.ReadFile(publicKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("公钥文件读取失败: %w", err)
	}
	block, _ := pem.Decode(fileBytes)
	if block == nil {
		return "", errors.New("公钥Decode错误")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("x509.ParsePKIXPublicKey：%w", err)
	}
	publicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("公钥解析错误")
	}
	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(originData))
	if err != nil {
		return "", fmt.Errorf("xrsa.EncryptPKCS1v15：%w", err)
	}
	return string(cipherBytes), nil
}
