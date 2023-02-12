# SWT
**Simple Web Token** inspired by [JWT](https://jwt.io/introduction), it can be used to validate data.

- [SWT](#swt)
	- [Background](#background)
	- [Install](#install)
	- [Usage](#usage)
	- [Examples](#examples)
	- [Licences](#licences)
## Background
While developing my personal blog, I encountered some scenarios that required authentication, so I used JWT, but for my small projects, JWT always produced long Token strings, so I designed SWT to create shorter Token strings.

Token structure
```
|--- 16 bits ---|--------|-- 128 bits --|
|- data length -|- data -|- sinagature -|
```

Signature: Use [MD5](https://pkg.go.dev/crypto/md5) to digest data, use [AES](https://pkg.go.dev/crypto/aes) to encrypt data
Use [base64](https://pkg.go.dev/encoding/base64) to encode Token.
## Install
```
go get github.com/hbread00/swt
```
## Usage
Create a Swt instance.
```
s, err := NewSwt([]byte("0123456789abcdef"))
if err != nil {
	panic(err)
}
```
Create a Token for your data.
```
token, err := s.MakeToken([]byte("your data"))
if err != nil {
	panic(err)
}
```
Verify token.
```
err := s.VerifyToken(token)
if err != nil {
	panic(err)
}
```
Parse data from token.
```
data, err := s.ParseData(token)
if err != nil {
	panic(err)
}
```
## Examples
```
func main() {
    s, err := NewSwt([]byte("0123456789abcdef"))
	if err != nil {
		panic(err)
	}
	data := []byte("sid: 4396, exp: 2200")
	fmt.Println("original data:", string(data))
	token, err := s.MakeToken(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("token:", token)
	err = s.VerifyToken(token)
	if err != nil {
		panic(err)
	}
	token_data, err := s.ParseData(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("data from token:", string(token_data))
}
```

## Licences
[MIT](LICENSE)