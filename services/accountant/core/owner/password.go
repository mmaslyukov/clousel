package owner

import (
	"crypto/sha256"
	"encoding/base64"
)

type IPassword interface {
	Hash() IPassword
	Base64() IPassword
	Str() string
	Eq(pswd IPassword) bool
}

type PasswordHashed struct {
	data string
}

func PasswordHashedCreate(data string) *PasswordHashed {
	return &PasswordHashed{data: data}
}

func (p *PasswordHashed) Hash() IPassword {
	return p
}

func (p *PasswordHashed) Base64() IPassword {
	return p
}

func (p *PasswordHashed) Str() string {
	return p.data
}

func (p *PasswordHashed) Eq(pswd IPassword) bool {
	return p.data == pswd.Str()
}

type PasswordPlain struct {
	data string
	hash []byte
	// hashed bool
}

func PasswordPlainCreate(data string) *PasswordPlain {
	return &PasswordPlain{data: data}
}

func (p *PasswordPlain) Hash() IPassword {
	if len(p.hash) == 0 {
		h := sha256.New()
		h.Write([]byte(p.data))
		p.hash = h.Sum(nil)
	}
	return p
}

func (p *PasswordPlain) Base64() IPassword {
	p.data = base64.StdEncoding.EncodeToString(p.hash)
	return p
}

func (p *PasswordPlain) Str() string {
	return p.data
}

func (p *PasswordPlain) Eq(pswd IPassword) bool {
	return p.data == pswd.Str()
}
