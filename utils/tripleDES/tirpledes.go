/*
 * @Author: Casso
 * @Date: 2021-02-02 09:52:23
 * @LastEditTime: 2021-02-02 09:55:46
 * @LastEditors: Please set LastEditors
 * @Description: 3DES
 * @FilePath: /githubStarChat/starChat/utils/tripleDES/tirpledes.go
 */

package tripledes

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

var GlobalTripleDES TripleDES

// TripleDES  key & iv
type TripleDES struct {
	Key string
	Iv  string
}

// Encrypt encrypt
func (t *TripleDES) Encrypt(plain string) (string, error) {
	key := []byte(t.Key)
	iv := []byte(t.Iv)

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	input := []byte(plain)
	// input = PKCS5Padding(input, block.BlockSize())
	blockMode := cipher.NewOFB(block, iv)
	crypted := make([]byte, len(input))
	blockMode.XORKeyStream(crypted, input)

	return base64.StdEncoding.EncodeToString(crypted), err
}

// Decrypt decry
func (t *TripleDES) Decrypt(secret string) (string, error) {
	key := []byte(t.Key)
	iv := []byte(t.Iv)

	crypted, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewOFB(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.XORKeyStream(origData, crypted)
	// origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

// PKCS5Padding padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding unpadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
