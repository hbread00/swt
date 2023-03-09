# SWT
**Simple Web Token** inspired by [JWT](https://jwt.io/introduction), it can be used to validate data. 

- [SWT](#swt)
	- [Background](#background)
	- [Install](#install)
	- [Usage](#usage)
	- [Examples](#examples)
	- [Licences](#licences)
## Background
I encountered some scenarios that required authentication. At first I used JWT, but the JWT Token string was too long for my small project, so I designed SWT to create a shorter Token string. 

Token structure
```
|-- 256 bit --|--------|
|- signature -|- data -|
```
Use [HMAC](https://pkg.go.dev/crypto/hmac)-[SHA256](https://pkg.go.dev/crypto/sha256) to digest and sign data.  
Use [Base64](https://pkg.go.dev/encoding/base64) to encode Token. 
## Install
```
go get github.com/hbread00/swt
```
## Usage
Create a Swt instance. 
```go
s := NewSwt([]byte("password"))
```
Create a Token for your data. 
```go
token, err := s.MakeToken([]byte("your data"))
if err != nil {
	panic(err)
}
```
Verify token. 
```go
err := s.VerifyToken(token)
if err != nil {
	panic(err)
}
```
Parse data from token. 
```go
data, err := s.ParseData(token)
if err != nil {
	panic(err)
}
```
Modify the key, no key is 100% infallible, the most secure key is the frequently modified key. 
```go
s.ResetSwt([]byte("password"))
```
## Examples
```go
func main() {
	s := NewSwt([]byte("0"))
	s.ResetSwt([]byte("password"))
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
Output:
```
original data: sid: 4396, exp: 2200
token: MSbLkwhP6fOvVRQEEIRfLO0DLU-ELM0uTJI3Ze7aEYNzaWQ6IDQzOTYsIGV4cDogMjIwMA
data from token: sid: 4396, exp: 2200
```
## Licences
[MIT](LICENSE)
