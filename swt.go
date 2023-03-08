package swt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
)

const (
	sIGN_LENGTH      = 32
	tOKEN_LENGTH_MIN = 44
)

var b64Encode = base64.RawURLEncoding

// Swt contains the encryption algorithm used
type Swt struct {
	key [64]byte
}

// Return a new Swt instance, encrypted key is required
func NewSwt(key []byte) *Swt {
	hash := sha512.Sum512(key)
	result := &Swt{
		key: hash,
	}
	return result
}

// Reset a Swt instance, encrypted key is required
func (s *Swt) ResetSwt(key []byte) {
	hash := sha512.Sum512(key)
	s.key = hash
}

// Enter data to create a Token
func (s *Swt) MakeToken(data []byte) (string, error) {
	if len(data) <= 0 {
		return "", errors.New("data length cannot be 0")
	}
	sign := s.sign(data)
	info := append(sign, data...)
	result := b64Encode.EncodeToString(info)
	return result, nil
}

// Verify the validity of the Token
// Need to use the same Swt as for encryption
func (s *Swt) VerifyToken(token string) error {
	if len(token) < tOKEN_LENGTH_MIN {
		return errors.New("invalid token length")
	}
	info, err := b64Encode.DecodeString(token)
	if err != nil {
		return err
	}
	data := info[sIGN_LENGTH:]
	sign := s.sign(data)
	if !hmac.Equal(sign, info[:sIGN_LENGTH]) {
		return errors.New("mismatch singature")
	}
	return nil
}

// Extracting data from Token
// This operation will not fully verify the legitimacy, please parse the Token after verifying its validity
func (s *Swt) ParseData(token string) ([]byte, error) {
	if len(token) < tOKEN_LENGTH_MIN {
		return nil, errors.New("invalid token length")
	}
	info, err := b64Encode.DecodeString(token)
	if err != nil {
		return nil, err
	}
	data := info[sIGN_LENGTH:]
	return data, nil
}

// Sign the data
func (s *Swt) sign(data []byte) []byte {
	mac := hmac.New(sha256.New, s.key[:])
	mac.Write(data)
	result := mac.Sum(nil)
	return result
}
