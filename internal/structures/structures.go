package structures

import (
	"crypto/sha256"
	"fmt"
)

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type HashedPassword struct {
	Hash string
	Salt string
}

func NewHashedPassword(password string, salt string) HashedPassword {
	hasher := sha256.New()
	fmt.Fprint(hasher, password, salt)

	return HashedPassword{
		Hash: string(hasher.Sum(nil)),
		Salt: salt,
	}
}

func (hp HashedPassword) Compare(password string) bool {
	hasher := sha256.New()
	fmt.Fprint(hasher, password, hp.Salt)

	return hp.Hash == string(hasher.Sum(nil))
}
