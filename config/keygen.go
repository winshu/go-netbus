package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"strings"
	"time"
)

// 固定格式
const timeLayout = "2006-01-02"

func padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func unPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//AES加密
func encrypt(original, key string) (string, error) {
	originalBytes := []byte(original)
	keyBytes := []byte(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	originalBytes = padding(originalBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	encrypted := make([]byte, len(originalBytes))
	blockMode.CryptBlocks(encrypted, originalBytes)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

//AES解密
func decrypt(encrypted, key string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyBytes[:blockSize])
	origData := make([]byte, len(encryptedBytes))
	blockMode.CryptBlocks(origData, encryptedBytes)
	origData = unPadding(origData)

	return string(origData), nil
}

func fixLength(seed string) string {
	if len(seed)%16 == 0 {
		return seed
	}
	return seed + strings.Repeat("=", 16-len(seed)%16)
}

// 生成 key
func NewKey(seed string, expired string) (string, error) {
	ex, err := time.Parse(timeLayout, expired)
	if err != nil {
		ex = time.Now().Add(30 * 24 * time.Hour)
	}
	expired = ex.Format(timeLayout)

	// 种子必须是 16的位数
	seed = fixLength(seed)
	return encrypt(expired, seed)
}

// 检查 key 是否有效
func CheckKey(seed, key string) (time.Time, bool) {
	// 超级 key
	if key == seed {
		return time.Time{}, true
	}

	seed = fixLength(seed)
	expired, err := decrypt(key, seed)
	if err != nil {
		log.Println("Fail to decrypt key.", err)
		return time.Time{}, false
	}

	ex, err := time.Parse(timeLayout, expired)
	if err != nil {
		log.Println("Fail to parse key", err)
	}
	return ex, time.Now().Before(ex)
}
