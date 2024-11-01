package utils

import (
	"github.com/Eggi19/simple-social-media/config"
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(pwd string) ([]byte, error)
	CheckPassword(pwd string, hash []byte) (bool, error)
}

type BCryptHasher struct {
}

func NewBCryptHasher() *BCryptHasher {
	return &BCryptHasher{}
}

func (b *BCryptHasher) HashPassword(pwd string) ([]byte, error) {
	con, err := config.ConfigInit()
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), con.HashCost)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (b *BCryptHasher) CheckPassword(pwd string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(pwd))
	if err != nil {
		return false, err
	}

	return true, nil
}