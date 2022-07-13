package main

import (
	"crypto/rand"
)

const (
	PBKDF2_ITER = 65536
	AES_KEYSIZE = 32
	AES_IV_SIZE = 16
	AES_NONCE_SIZE = 16
	AES_GCM_AAD = "A256GCM"
)

func GenCryptoRand(nb uint8) []byte {
	data := make([]byte, nb)
	_, err := rand.Read(data)
	if err != nil {
		panic("Cannot generate random bytes !")
	}
	return data
}
