package swt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
)

// Swt contains the encryption algorithm used
type Swt struct {
	encrypter cipher.Block
	salt      []byte
}

// Return a new Swt instance, encrypted key is required
// Since Swt uses Advanced Encryption Standard, the key length should conform to the AES specification(128/192/256bits)
func NewSwt(key []byte) (*Swt, error) {
	aes_encrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	salt := sha512.Sum512_256(key)
	result := &Swt{
		encrypter: aes_encrypter,
		salt:      salt[8:24],
	}
	return result, nil
}

// Reset a Swt instance, encrypted key is required
// The key length should conform to the AES specification(128/192/256bits)
func (s *Swt) ResetSwt(key []byte) error {
	aes_encrypter, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	salt := sha512.Sum512_256(key)
	s.encrypter = aes_encrypter
	s.salt = salt[8:24]
	return nil
}

// Enter data to create a Token
func (s *Swt) MakeToken(data []byte) (string, error) {
	if len(data) <= 0 {
		return "", errors.New("invalid data length")
	}
	sign := s.sign(data)
	info := append(sign, data...)
	result := base64.StdEncoding.EncodeToString(info)
	return result, nil
}

// Verify the validity of the Token
// Need to use the same Swt as for encryption
func (s *Swt) VerifyToken(token string) error {
	if len(token) < 24 {
		return errors.New("invalid token length")
	}
	info, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}
	data := info[16:]
	sign := s.sign(data)
	if !compare(sign, info[:16]) {
		return errors.New("invalid singature")
	}
	return nil
}

// Extracting data from Token
// This operation will not fully verify the legitimacy, please parse the Token after verifying its validity
func (s *Swt) ParseData(token string) ([]byte, error) {
	if len(token) < 24 {
		return nil, errors.New("invalid token length")
	}
	info, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	data := info[16:]
	return data, nil
}

// Sign the data
func (s *Swt) sign(data []byte) []byte {
	info := append(data, s.salt...)
	hash := sha256.Sum256(info)
	result := make([]byte, 16)
	s.encrypter.Encrypt(result, hash[8:24])
	return result
}

// Compare two slices
func compare(lhs []byte, rhs []byte) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for i := range lhs {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}
