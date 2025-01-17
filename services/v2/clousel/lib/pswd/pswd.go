package pswd

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

type IPassword interface {
	Str() string
}

type IPasswordPlain interface {
	Str() string
	Hash() IPasswordHashed
}

type IPasswordHashed interface {
	HexStr() string
	Encode() IPasswordHasedBase64
}

type IPasswordHasedBase64 interface {
	Data() *[]byte
	HexStr() string
	Str() string
	Decode() (IPasswordHashed, error)
	Equal(p IPasswordHasedBase64) bool
}

type Password struct {
	data []byte
}

func (p *Password) Str() string {
	return string(p.data[:])
}

func (p *Password) HexStr() string {
	return hex.EncodeToString(p.data)
}

func PasswordPlainCreate(pswd string) IPasswordPlain {
	return &Password{data: []byte(pswd)}
}
func PasswordHasedBase64Create(pswd string) IPasswordHasedBase64 {
	return &Password{data: []byte(pswd)}
}

func (p *Password) Hash() IPasswordHashed {
	// return &p
	var n Password
	h := sha256.New()
	h.Write(p.data)
	n.data = h.Sum(nil)
	return &n
}

func (p *Password) Encode() IPasswordHasedBase64 {
	var n Password
	n.data = make([]byte, base64.StdEncoding.EncodedLen(len(p.data)))
	base64.StdEncoding.Encode(n.data, p.data)
	return &n
}

func (p *Password) Decode() (IPasswordHashed, error) {
	var n Password
	n.data = make([]byte, base64.StdEncoding.DecodedLen(len(p.data)))
	l, e := base64.StdEncoding.Decode(n.data, p.data)
	if e == nil {
		n.data = n.data[:l]
	}
	return &n, e
}

func (p *Password) Data() *[]byte {
	return &p.data
}

func (p *Password) Equal(po IPasswordHasedBase64) bool {
	return bytes.Equal(p.data, *po.Data())
}
