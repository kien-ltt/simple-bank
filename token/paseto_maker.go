package token

import (
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	pasetoSymmetricKey, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		symmetricKey: pasetoSymmetricKey,
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token := paseto.NewToken()
	token.Set("id", payload.ID)
	token.Set("username", payload.Username)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiresAt)

	return token.V4Encrypt(maker.symmetricKey, nil), payload, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, nil)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, err := getPayloadFromPasetoToken(parsedToken)
	return payload, err
}

func getPayloadFromPasetoToken(t *paseto.Token) (*Payload, error) {
	id, err := t.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}
	username, err := t.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}
	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiresAt: expiredAt,
	}, nil
}
