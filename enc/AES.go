package enc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// aes加密
func Encrypt(plantText, key []byte) ([]byte, error) {
	// 这个代码比较重要，是涉及到签名大小
	if len(key) > 16 {
		key = key[:16]
	}
	//NewCipher该函数限制了输入k的长度必须为16, 24或者32l
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//补全码
	plantText = PKCS7Padding(plantText, block.BlockSize())
	//加密模式
	blockModel := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	//创建数组
	ciphertext := make([]byte, len(plantText))
	//加密
	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// aes解密
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	// 这个代码比较重要，是涉及到签名大小
	if len(key) > 16 {
		key = key[:16]
	}
	//分组秘钥
	block, err := aes.NewCipher(key) // 选择加密算法
	if err != nil {
		return nil, err
	}
	//加密模式
	blockModel := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	plantText := make([]byte, len(ciphertext))
	//解密
	blockModel.CryptBlocks(plantText, ciphertext)
	//去补全码
	plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
