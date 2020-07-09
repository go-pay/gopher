package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
)

// 解密数据的Bytes数组
func AesDecryptToBytes(data, secretKey string) ([]byte, error) {
	secretData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	key := []byte(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(secretData))
	blockMode.CryptBlocks(originData, secretData)
	originData = PKCS5UnPadding(originData)
	return originData, nil
}

// 解密数据到结构体
func AesDecryptToStruct(data, secretKey string, beanPtr interface{}) (err error) {
	//验证参数类型
	beanValue := reflect.ValueOf(beanPtr)
	if beanValue.Kind() != reflect.Ptr {
		return errors.New("传入参数类型必须是以指针形式")
	}
	//验证interface{}类型
	if beanValue.Elem().Kind() != reflect.Struct {
		return errors.New("传入interface{}必须是结构体")
	}
	secretData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	key := []byte(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(secretData))
	blockMode.CryptBlocks(originData, secretData)
	originData = PKCS5UnPadding(originData)
	//解析
	err = json.Unmarshal(originData, beanPtr)
	if err != nil {
		return err
	}
	return nil
}

// 解密数据到Map集合
func AesDecryptToMap(data, secretKey string) (mapData map[string]interface{}, err error) {
	secretData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	key := []byte(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(secretData))
	blockMode.CryptBlocks(originData, secretData)
	originData = PKCS5UnPadding(originData)
	//解析
	mapData = make(map[string]interface{}, 0)
	err = json.Unmarshal(originData, &mapData)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}
