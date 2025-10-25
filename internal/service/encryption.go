package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

// TokenData represents the data structure to be encrypted in the token
type TokenData struct {
	SurveyID    uint                   `json:"survey_id"`
	PrefillData map[string]interface{} `json:"prefill_data"`
	ExpiresAt   int64                  `json:"expires_at"`
	UniqueID    string                 `json:"unique_id"`
}

// EncryptionService defines the interface for encryption operations
type EncryptionService interface {
	EncryptToken(data *TokenData) (string, error)
	DecryptToken(token string) (*TokenData, error)
}

// encryptionService implements EncryptionService using AES-256-GCM
type encryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service instance
// key must be exactly 32 bytes for AES-256
func NewEncryptionService(key string) (EncryptionService, error) {
	keyBytes := []byte(key)
	
	// Validate key length
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("encryption key must be exactly 32 bytes, got %d bytes", len(keyBytes))
	}
	
	return &encryptionService{
		key: keyBytes,
	}, nil
}

// EncryptToken encrypts TokenData and returns a base64 URL-safe encoded string
func (s *encryptionService) EncryptToken(data *TokenData) (string, error) {
	// Serialize TokenData to JSON
	plaintext, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token data: %w", err)
	}
	
	// Create AES cipher block
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Generate random nonce (IV)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	// Encrypt the plaintext
	// The nonce is prepended to the ciphertext
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	
	// Encode to base64 URL-safe format
	encoded := base64.URLEncoding.EncodeToString(ciphertext)
	
	return encoded, nil
}

// DecryptToken decrypts a base64 URL-safe encoded token and returns TokenData
func (s *encryptionService) DecryptToken(token string) (*TokenData, error) {
	// Decode from base64 URL-safe format
	ciphertext, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}
	
	// Create AES cipher block
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Validate ciphertext length
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	
	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt token: %w", err)
	}
	
	// Deserialize JSON to TokenData
	var data TokenData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}
	
	return &data, nil
}
