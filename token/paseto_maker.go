package token

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be %d characters", chacha20poly1305.KeySize)
	}

	v4SymKey, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))

	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		symmetricKey: v4SymKey,
	}

	return maker, nil
}

func (p PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	token := paseto.NewToken()

	token.SetIssuer("simplebank")
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	token.SetString("id", uuid.NewString())
	token.SetString("username", username)

	return token.V4Encrypt(p.symmetricKey, nil), nil

}

func (p PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()

	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.IssuedBy("simplebank"))
	parser.AddRule(paseto.ValidAt(time.Now()))

	decryptedToken, err := parser.ParseV4Local(p.symmetricKey, token, nil)

	if err != nil {
		return nil, err
	}

	username, err := decryptedToken.GetString("username")
	if err != nil {
		return nil, err
	}
	id, err := decryptedToken.GetString("id")

	if err != nil {
		return nil, err
	}

	issuedAt, err := decryptedToken.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	expiredAt, err := decryptedToken.GetExpiration()

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &Payload{Username: username, ID: id, IssuedAt: issuedAt, ExpiredAt: expiredAt}, nil
}

//func ForAudience(audience string) Rule
//func IdentifiedBy(identifier string) Rule
//func IssuedBy(issuer string) Rule
//func NotExpired() Rule
//func Subject(subject string) Rule
//func ValidAt(t time.Time) Rule
