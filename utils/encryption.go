package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// padding 使用PKCS7进行填充
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// pkcs7UnPadding 使用PKCS7进行去填充
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// EncryptAES128CBC 使用AES-128 CBC模式进行加密
func EncryptAES128CBC(plainText []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	// 使用PKCS7进行填充
	plainText = pkcs7Padding(plainText, blockSize)
	ciphertext := make([]byte, len(plainText))
	mode := cipher.NewCBCEncrypter(block, key[:blockSize])
	mode.CryptBlocks(ciphertext, plainText)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES128CBC 使用AES-128 CBC模式进行解密
func DecryptAES128CBC(cipherText string, key []byte) ([]byte, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	ciphertext := make([]byte, len(cipherTextBytes))
	mode := cipher.NewCBCDecrypter(block, key[:blockSize])
	mode.CryptBlocks(ciphertext, cipherTextBytes)
	// 去除PKCS7填充
	ciphertext = pkcs7UnPadding(ciphertext)
	return ciphertext, nil
}
