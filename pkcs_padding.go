//==================================
//  * Name：Jerry
//  * DateTime：2019/5/17 13:06
//  * Desc：
//==================================
package go_util

import "bytes"

//加密填充模式（添加补全码） PKCS5Padding
//加密时，如果加密bytes的length不是blockSize的整数倍，需要在最后面添加填充byte
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	paddingCount := blockSize - len(ciphertext)%blockSize //需要padding的数目
	paddingBytes := []byte{byte(paddingCount)}
	padtext := bytes.Repeat(paddingBytes, paddingCount) //生成填充的文本
	return append(ciphertext, padtext...)
}

//解密填充模式（去除补全码） PKCS5UnPadding
//解密时，需要在最后面去掉加密时添加的填充byte
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])   //找到Byte数组最后的填充byte
	return origData[:(length - unpadding)] //只截取返回有效数字内的byte数组
}
