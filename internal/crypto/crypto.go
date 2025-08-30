package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"

	"doc0x1/text2babe/internal/config"
)

func EncryptData(data string, cfg *config.Config) (string, error) {
	// Always treat input as text
	inputBytes := []byte(data)
	
	var outputBytes []byte
	
	if cfg.UseEncryption {
		// AES-GCM encryption
		block, err := aes.NewCipher(cfg.Key)
		if err != nil {
			return "", fmt.Errorf("failed to create cipher: %w", err)
		}
		
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", fmt.Errorf("failed to create GCM: %w", err)
		}
		
		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return "", fmt.Errorf("failed to generate nonce: %w", err)
		}
		
		outputBytes = gcm.Seal(nonce, nonce, inputBytes, nil)
	} else {
		// Plain encoding - just use the input bytes directly
		outputBytes = inputBytes
	}
	
	// Apply output format
	switch cfg.OutputMode {
	case "hex":
		return hex.EncodeToString(outputBytes), nil
	case "base64":
		return base64.StdEncoding.EncodeToString(outputBytes), nil
	case "binary":
		// For binary output, we need to ensure it's displayable
		// Convert to a readable binary representation (0s and 1s)
		var binaryStr string
		for _, b := range outputBytes {
			binaryStr += fmt.Sprintf("%08b", b)
		}
		return binaryStr, nil
	default:
		return hex.EncodeToString(outputBytes), nil
	}
}

func DecryptData(data string, cfg *config.Config) (string, error) {
	var inputBytes []byte
	var err error
	
	// Auto-detect input format with smart binary detection
	// Check if data looks like binary (only 0s and 1s, length divisible by 8)
	isBinaryFormat := len(data) > 8 && len(data)%8 == 0
	for _, r := range data {
		if r != '0' && r != '1' {
			isBinaryFormat = false
			break
		}
	}
	
	if isBinaryFormat {
		// If it looks like binary format, try binary parsing first
		inputBytes, err = parseBinaryString(data)
		if err != nil {
			// If binary parsing fails, fall back to other formats
			inputBytes, err = hex.DecodeString(data)
			if err != nil {
				inputBytes, err = base64.StdEncoding.DecodeString(data)
				if err != nil {
					inputBytes = []byte(data)
				}
			}
		}
	} else {
		// For non-binary data, try hex first, then base64, then raw binary
		inputBytes, err = hex.DecodeString(data)
		if err != nil {
			// If hex fails, try base64
			inputBytes, err = base64.StdEncoding.DecodeString(data)
			if err != nil {
				// If base64 fails, treat as raw binary
				inputBytes = []byte(data)
			}
		}
	}
	
	if cfg.UseEncryption {
		// AES-GCM decryption
		block, err := aes.NewCipher(cfg.Key)
		if err != nil {
			return "", fmt.Errorf("failed to create cipher: %w", err)
		}
		
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", fmt.Errorf("failed to create GCM: %w", err)
		}
		
		nonceSize := gcm.NonceSize()
		if len(inputBytes) < nonceSize {
			return "", fmt.Errorf("ciphertext too short")
		}
		
		nonce, ciphertext := inputBytes[:nonceSize], inputBytes[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt data: %w", err)
		}
		
		return string(plaintext), nil
	} else {
		// Plain decoding - just convert back to string
		return string(inputBytes), nil
	}
}

// parseBinaryString converts a string of 0s and 1s to bytes
func parseBinaryString(binaryStr string) ([]byte, error) {
	// Remove any whitespace and newlines
	binaryStr = strings.ReplaceAll(binaryStr, " ", "")
	binaryStr = strings.ReplaceAll(binaryStr, "\n", "")
	binaryStr = strings.ReplaceAll(binaryStr, "\r", "")
	binaryStr = strings.ReplaceAll(binaryStr, "\t", "")
	
	// Check if the string is long enough and contains only 0s and 1s
	if len(binaryStr) == 0 {
		return nil, fmt.Errorf("empty binary string")
	}
	
	// Check if it's a valid binary string (only 0s and 1s)
	validChars := 0
	for _, r := range binaryStr {
		if r == '0' || r == '1' {
			validChars++
		} else {
			return nil, fmt.Errorf("invalid binary character: %c at position %d", r, validChars)
		}
	}
	
	// Binary string must be divisible by 8 to form complete bytes
	if len(binaryStr)%8 != 0 {
		return nil, fmt.Errorf("binary string length (%d) must be divisible by 8", len(binaryStr))
	}
	
	var result []byte
	for i := 0; i < len(binaryStr); i += 8 {
		byteStr := binaryStr[i : i+8]
		b, err := strconv.ParseUint(byteStr, 2, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to parse binary byte '%s' at position %d: %w", byteStr, i/8, err)
		}
		result = append(result, byte(b))
	}
	
	return result, nil
}