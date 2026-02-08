package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type Crypto struct {
	key []byte
}

func NewCrypto(secretKey string) (*Crypto, error) {
	if len(secretKey) != 32 {
		return nil, fmt.Errorf("secret key muss genau 32 Zeichen (256-bit) lang sein, ist aber %d Zeichen", len(secretKey))
	}
	return &Crypto{
		key: []byte(secretKey),
	}, nil
}

func (c *Crypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("fehler beim Erstellen des AES-Cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("fehler beim Erstellen von GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("fehler beim Generieren der Nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Crypto) Decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("fehler beim Dekodieren von Base64: %w", err)
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("fehler beim Erstellen des AES-Cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("fehler beim Erstellen von GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext zu kurz")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("fehler beim Entschlüsseln: %w", err)
	}

	return string(plaintext), nil
}

func (c *Crypto) MustEncrypt(plaintext string) string {
	encrypted, err := c.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}
	return encrypted
}

func (c *Crypto) MustDecrypt(encoded string) string {
	decrypted, err := c.Decrypt(encoded)
	if err != nil {
		panic(err)
	}
	return decrypted
}