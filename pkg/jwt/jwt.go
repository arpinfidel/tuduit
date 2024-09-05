package jwt

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SigningMethod string

	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New(signingMethod string, privateKey, publicKey []byte) (*JWT, error) {
	j := &JWT{
		SigningMethod: signingMethod,
	}

	priv, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}
	j.privateKey = priv

	pub, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}
	j.publicKey = pub

	return j, nil
}

func (j *JWT) Sign(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.GetSigningMethod(j.SigningMethod), claims).SignedString(j.privateKey)
}

func (j *JWT) Verify(token string, targetClaims jwt.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, targetClaims, func(token *jwt.Token) (interface{}, error) {
		return j.publicKey, nil
	})
}
