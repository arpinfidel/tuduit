package app

import (
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/golang-jwt/jwt/v5"
)

func (a *App) makeAccessToken(u entity.User) (string, error) {
	accessTok, err := a.d.JWT.Sign(entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
		},
		TokenType: "access",
		UserID:    u.ID,
	})
	if err != nil {
		return "", err
	}

	return accessTok, nil
}

func (a *App) makeRefreshToken(u entity.User) (string, error) {
	refreshTok, err := a.d.JWT.Sign(entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		TokenType: "refresh",
		UserID:    u.ID,
	})
	if err != nil {
		return "", err
	}

	return refreshTok, nil
}

func (a *App) makeTokenPair(u entity.User) (accessToken, refreshToken string, err error) {
	accessTok, err := a.makeAccessToken(u)
	if err != nil {
		return "", "", err
	}

	refreshTok, err := a.makeRefreshToken(u)
	if err != nil {
		return "", "", err
	}

	return accessTok, refreshTok, nil
}

func (a *App) VerifyToken(token string) (entity.Claims, error) {
	var claims entity.Claims

	_, err := a.d.JWT.Verify(token, &claims)
	if err != nil {
		return entity.Claims{}, err
	}

	return claims, nil
}
