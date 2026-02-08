package security

import (
	"testing"
)

func TestNewCrypto(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		wantErr   bool
	}{
		{
			name:      "valid 32 character key",
			secretKey: "12345678901234567890123456789012",
			wantErr:   false,
		},
		{
			name:      "key too short",
			secretKey: "short",
			wantErr:   true,
		},
		{
			name:      "key too long",
			secretKey: "123456789012345678901234567890123",
			wantErr:   true,
		},
		{
			name:      "empty key",
			secretKey: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCrypto(tt.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCrypto() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCrypto_EncryptDecrypt(t *testing.T) {
	crypto, err := NewCrypto("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("Failed to create crypto: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "Hello, World!",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "long text",
			plaintext: "This is a very long text that should still be encrypted and decrypted correctly without any issues. It contains multiple sentences and special characters like !@#$%^&*().",
		},
		{
			name:      "unicode characters",
			plaintext: "Hello 世界 🌍",
		},
		{
			name:      "json-like string",
			plaintext: `{"key": "value", "number": 123, "boolean": true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := crypto.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			// Decrypt
			decrypted, err := crypto.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			// Verify
			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.plaintext)
			}

			// Verify that encrypted value is different from plaintext
			if tt.plaintext != "" && encrypted == tt.plaintext {
				t.Error("Encrypted value should be different from plaintext")
			}
		})
	}
}

func TestCrypto_EncryptProducesDifferentCiphertext(t *testing.T) {
	crypto, err := NewCrypto("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("Failed to create crypto: %v", err)
	}

	plaintext := "Same plaintext"

	// Encrypt the same plaintext multiple times
	encrypted1, err := crypto.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	encrypted2, err := crypto.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Due to random nonce, ciphertexts should be different
	if encrypted1 == encrypted2 {
		t.Error("Multiple encryptions of the same plaintext should produce different ciphertexts")
	}

	// But both should decrypt to the same plaintext
	decrypted1, _ := crypto.Decrypt(encrypted1)
	decrypted2, _ := crypto.Decrypt(encrypted2)

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Error("Both encrypted values should decrypt to the original plaintext")
	}
}

func TestCrypto_DecryptInvalidData(t *testing.T) {
	crypto, err := NewCrypto("12345678901234567890123456789012")
	if err != nil {
		t.Fatalf("Failed to create crypto: %v", err)
	}

	tests := []struct {
		name      string
		encrypted string
	}{
		{
			name:      "invalid base64",
			encrypted: "not-valid-base64!@#$",
		},
		{
			name:      "too short ciphertext",
			encrypted: "dGVzdA==", // "test" in base64, but too short
		},
		{
			name:      "corrupted data",
			encrypted: "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := crypto.Decrypt(tt.encrypted)
			if err == nil {
				t.Error("Decrypt() should return an error for invalid data")
			}
		})
	}
}

func BenchmarkCrypto_Encrypt(b *testing.B) {
	crypto, _ := NewCrypto("12345678901234567890123456789012")
	plaintext := "This is a test plaintext for benchmarking encryption performance"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = crypto.Encrypt(plaintext)
	}
}

func BenchmarkCrypto_Decrypt(b *testing.B) {
	crypto, _ := NewCrypto("12345678901234567890123456789012")
	plaintext := "This is a test plaintext for benchmarking decryption performance"
	encrypted, _ := crypto.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = crypto.Decrypt(encrypted)
	}
}