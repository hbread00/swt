package swt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"errors"
)

// Swt contains the encryption algorithm used
type Swt struct {
	encrypter cipher.Block
}

// Return a new Swt instance, encrypted key is required
// Since Swt uses Advanced Encryption Standard, the key length should conform to the AES specification(128/192/256bits)
func NewSwt(key []byte) (*Swt, error) {
	aes_encrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	result := &Swt{
		encrypter: aes_encrypter,
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
	s.encrypter = aes_encrypter
	return nil
}

// Enter data to create a Token
func (s *Swt) MakeToken(data []byte) (string, error) {
	if len(data) <= 0 || len(data) > int(^uint16(0)) {
		return "", errors.New("invalid data length")
	}
	info := make([]byte, 2)
	binary.BigEndian.PutUint16(info, uint16(len(data)))
	info = append(info, data...)
	sign := s.sign(info)
	info = append(info, sign...)
	result := base64.StdEncoding.EncodeToString(info)
	return result, nil
}

// Verify the validity of the Token
// Need to use the same Swt as for encryption
func (s *Swt) VerifyToken(token string) error {
	info, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}
	if len(info) < 16+2 {
		return errors.New("invalid token length")
	}
	data_len := binary.BigEndian.Uint16(info[0:2])
	if len(info) != int(data_len)+2+16 {
		return errors.New("invalid token structure")
	}
	data := info[:data_len+2]
	sign := s.sign(data)
	if !compare(sign, info[len(info)-16:]) {
		return errors.New("invalid singature")
	}
	return nil
}

// Extracting data from Token
// This operation will not fully verify the legitimacy, please parse the Token after verifying its validity
func (s *Swt) ParseData(token string) ([]byte, error) {
	info, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	data_len := binary.BigEndian.Uint16(info[0:2])
	if len(info) != int(data_len)+2+16 {
		return nil, errors.New("invalid token structure")
	}
	data := info[2 : data_len+2]
	return data, nil
}

// Sign the data
func (s *Swt) sign(data []byte) []byte {
	hash := md5.Sum(data)
	result := make([]byte, 16)
	s.encrypter.Encrypt(result, hash[:])
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
