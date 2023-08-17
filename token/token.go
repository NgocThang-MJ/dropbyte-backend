package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type Payload struct {
	TokenId   uuid.UUID `json:"token_id"`
	UserId    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Paseto struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

type Token interface {
	CreateToken(userId uuid.UUID, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

func NewMaker(key string) (Token, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf(
			"Invalid key size: must be exactly %d characters",
			chacha20poly1305.KeySize,
		)
	}
	maker := &Paseto{paseto: paseto.NewV2(), symmetricKey: []byte(key)}
	return maker, nil
}

func (maker *Paseto) CreateToken(userId uuid.UUID, duration time.Duration) (string, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	payload := &Payload{
		TokenId:   tokenId,
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, nil
}

func (maker *Paseto) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.CheckExpired()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (payload *Payload) CheckExpired() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("Token has expired")
	}
	return nil
}
