package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// paddingの追加
	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	if padding != 0 {
		paddingBytes := bytes.Repeat([]byte{byte(padding)}, padding)
		plaintext = append(plaintext, paddingBytes...)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// 初期化ベクトル（IV）をランダムに作成し、暗号文の先頭に追加
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// 暗号化
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// 暗号文から初期化ベクトル（IV）を取得し、IV以降のデータを復号化対象とする
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 復号化用のバッファを用意
	decrypted := make([]byte, len(ciphertext))

	// 復号化
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, ciphertext)

	// paddingの削除
	padding := int(decrypted[len(decrypted)-1])
	plaintext := decrypted[:len(decrypted)-padding]
	return plaintext, nil
}
