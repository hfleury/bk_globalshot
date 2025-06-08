package token

import (
	"fmt"
	"time"

	paseto "aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	privateKey   paseto.V4AsymmetricSecretKey
	PublicKeyHex paseto.V4AsymmetricPublicKey
}

func NewPasetoMaker(privateKeyHex string) (Maker, error) {
	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &PasetoMaker{
		privateKey:   privateKey,
		PublicKeyHex: privateKey.Public(),
	}, nil
}

func (m *PasetoMaker) CreateToken(email string, duration time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(duration)

	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(exp)
	token.SetString("email", email)

	signedToken := token.V4Sign(m.privateKey, nil)
	return signedToken, nil
}

func (m *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Public(m.PublicKeyHex, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	email, err := parsedToken.GetString("email")
	if err != nil {
		return nil, fmt.Errorf("missing 'email' claim")
	}

	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return nil, fmt.Errorf("get issued at error %w", err)
	}

	expiration, err := parsedToken.GetExpiration()
	if err != nil {
		return nil, fmt.Errorf("get expiration error %w", err)
	}

	payload := &Payload{
		Email:     email,
		IssuedAt:  issuedAt,
		ExpiresAt: expiration,
	}

	return payload, nil
}
