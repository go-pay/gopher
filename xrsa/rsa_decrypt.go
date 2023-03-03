package xrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
)

// RSA解密数据
// t：PKCS1 或 PKCS8
// cipherData：加密字符串byte数组
// privateKey：私钥
func RsaDecrypt(t PKCSType, cipherData []byte, privateKey string) (originData []byte, err error) {
	var (
		key *rsa.PrivateKey
	)

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("privateKey decode error")
	}

	switch t {
	case PKCS1:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	case PKCS8:
		pkcs8Key, e := x509.ParsePKCS8PrivateKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		pk8, ok := pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("parse PKCS8 key error")
		}
		key = pk8
	default:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	originBytes, err := rsa.DecryptPKCS1v15(rand.Reader, key, cipherData)
	if err != nil {
		return nil, fmt.Errorf("xrsa.DecryptPKCS1v15：%w", err)
	}
	return originBytes, nil
}

// RSA解密数据
// OAEPWithSHA-256AndMGF1Padding
func RsaDecryptOAEP(h hash.Hash, t PKCSType, privateKey string, ciphertext, label []byte) (originData []byte, err error) {
	var (
		key *rsa.PrivateKey
	)

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("privateKey decode error")
	}

	switch t {
	case PKCS1:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	case PKCS8:
		pkcs8Key, e := x509.ParsePKCS8PrivateKey(block.Bytes)
		if e != nil {
			return nil, e
		}
		pk8, ok := pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("parse PKCS8 key error")
		}
		key = pk8
	default:
		if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	originBytes, err := rsa.DecryptOAEP(h, rand.Reader, key, ciphertext, label)
	if err != nil {
		return nil, err
	}
	return originBytes, nil
}
