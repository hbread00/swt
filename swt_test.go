package swt

import (
	"fmt"
	"strconv"
	"testing"
)

func TestQucikStart(t *testing.T) {
	s, err := NewSwt([]byte("0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("sid: 4396, exp: 2200")
	fmt.Println("original data:", string(data))
	token, err := s.MakeToken(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("token:", token)
	err = s.VerifyToken(token)
	if err != nil {
		t.Fatal(err)
	}
	token_data, err := s.ParseData(token)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("data from token:", string(token_data))
}

func TestMakeToken(t *testing.T) {
	s, err := NewSwt([]byte("0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	var max_len_data [int(^uint16(0))]byte
	var too_long_data [int(^uint16(0)) + 1]byte
	cases := []struct {
		name   string
		input  []byte
		target bool
	}{
		{
			"normal",
			[]byte("test data"),
			true,
		},
		{
			"empty input",
			[]byte{},
			false,
		},
		{
			"short data",
			[]byte("0"),
			true,
		},
		{
			"max length data",
			(max_len_data)[:],
			true,
		},
		{
			"too long data",
			too_long_data[:],
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, result := s.MakeToken(c.input)
			if (result == nil) != c.target {
				t.Errorf("input: %d | target: %t | result: %t", c.input, c.target, (result == nil))
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	s, err := NewSwt([]byte("0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name   string
		input  string
		target bool
	}{
		{
			"normal",
			"AB50ZXN0IGRhdGEsIGlkOiB0ZXN0LCBleHA6IHRlc3QM4XH3+YqgnHNPxYGMcuPj",
			true,
		},
		{
			"not base64",
			"AB50ZXN0IGRhdGEsIGlkOiB0ZXN0LCBleHA6IHRlc3QM4XH3+YqgnHNPxYGMcuP",
			false,
		},
		{
			"wrong length",
			"AB50ZXN0IGRhdGEsIGlkOiB0ZXN0LCBleHA6IHRlc3QM4XH3+YqgnHNPxYGM",
			false,
		},
		{
			"wrong data length",
			"aB50ZXN0IGRhdGEsIGlkOiB0ZXN0LCBleHA6IHRlc3QM4XH3+YqgnHNPxYGM",
			false,
		},
		{
			"wrong data",
			"AB50ZXN0IGRhdgEsIGlkOiB0ZXN0LCBleHA6IHRlc3QM4XH3+YqgnHNPxYGMcuPj",
			false,
		},
		{
			"wrong signature",
			"AB50ZXN0IGRhdGEsIGlkOiB0ZXN0LCBleHA6IHRlc3QM4XH3+YqgnHNPxYgMcuPj",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := s.VerifyToken(c.input)
			fmt.Println("err:", result)
			if (result == nil) != c.target {
				t.Errorf("input: %s | target: %t | result: %t", c.input, c.target, (result == nil))
			}
		})
	}
}

func TestParseData(t *testing.T) {
	s, err := NewSwt([]byte("0123456789abcdef"))
	if err != nil {
		t.Fatal(err)
	}
	data_base := "test"
	for i := 0; i < 10; i += 1 {
		t.Run("test "+strconv.Itoa(i), func(t *testing.T) {
			input := []byte(data_base + strconv.Itoa(i) + data_base)
			token, _ := s.MakeToken(input)
			fmt.Println("token:", token)
			result, _ := s.ParseData(token)
			if !compare(result, input) {
				t.Errorf("input: %s | target: %s | result: %s", input, input, result)
			}
		})
	}
}
