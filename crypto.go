package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

const (
	PBKDF2_SALT_SIZE = 16
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

// Generate a derived key from a pawwsord.
// Return the derived key and the salt.
func GenPBKDF2(pwd []byte) ([]byte, []byte) {
	salt := GenCryptoRand(PBKDF2_SALT_SIZE)
	key := pbkdf2.Key(pwd, salt, PBKDF2_ITER, AES_KEYSIZE, sha256.New)
	return key, salt
}

// Generate a derived key from a pawwsord and salt.
// Return the derived key.
func GenPBKDF2WithSalt(pwd []byte, salt []byte) []byte {
	key := pbkdf2.Key(pwd, salt, PBKDF2_ITER, AES_KEYSIZE, sha256.New)
	return key
}

// Encrypt the 'plain' message with AES_256_GCM and the given 'key'.
// This function will generate a crypto random IV used for encryption.
// Return the encrypted message and the IV.
func EncryptWithKey(plain *[]byte, key *[]byte) (*[]byte, *[]byte, error) {
	iv := GenCryptoRand(AES_IV_SIZE)
	block, err := aes.NewCipher(*key)
	if err != nil {
		return &[]byte{}, &[]byte{}, err
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, AES_NONCE_SIZE)
	if err != nil {
		return &[]byte{}, &[]byte{}, err
	}
	cipher := gcm.Seal(nil, iv, *plain, []byte(AES_GCM_AAD))
	return &cipher, &iv, nil
}

// Decrypt some []byte with AES_256_GCM and the given key and iv.
func Decrypt(encrypted *[]byte, key *[]byte, iv *[]byte) (*[]byte, error) {
	block, err := aes.NewCipher(*key)
	if err != nil {
		return &[]byte{}, err
	}
	gcm, err := cipher.NewGCMWithNonceSize(block, len(*iv))
	if err != nil {
		return &[]byte{}, err
	}
	plain, err := gcm.Open(nil, *iv, *encrypted, []byte(AES_GCM_AAD))
	if err != nil {
		return &[]byte{}, err
	}
	return &plain, nil
}
