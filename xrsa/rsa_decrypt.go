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

// RSA解密数据
// cipherData：加密字符串
// privateKeyFilePath：私钥证书文件路径
func RsaDecryptData(cipherData string, privateKeyFilePath string) (originData string, err error) {
	fileBytes, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("私钥文件读取失败: %w", err)
	}
	block, _ := pem.Decode(fileBytes)
	if block == nil {
		return "", errors.New("私钥Decode错误")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("x509.ParsePKCS1PrivateKey：%w", err)
	}
	originBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(cipherData))
	if err != nil {
		return "", fmt.Errorf("xrsa.DecryptPKCS1v15：%w", err)
	}
	return string(originBytes), nil
}
