package utils

import (
	"crypto/rc4"
	"fmt"
)

func RC4Encrypt(bytes []byte) []byte {
	c, err := rc4.NewCipher([]byte(GlobalConfig.RC4Key))

	if err != nil {
		fmt.Println("Encrypt err", err)
		return bytes
	}

	dst := make([]byte, len(bytes))

	c.XORKeyStream(dst, bytes)

	return dst
}

func RC4Decrypt(bytes []byte) []byte {
	c, err := rc4.NewCipher([]byte(GlobalConfig.RC4Key))
	if err != nil {
		fmt.Println("Decrypt err", err)
		return bytes
	}

	dst := make([]byte, len(bytes))

	c.XORKeyStream(dst, bytes)

	return dst
}
