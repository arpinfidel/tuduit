package app

import (
	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/crypto"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/errs"
)

type UserLoginParams struct {
	WhatsAppNumber string `rose:"whatsapp_number,wa"`
	UserName       string `rose:"user_name,username"`
	Password       string `rose:"password,pw,required="`
}

type UserLoginResponse struct {
	AccessToken  string `rose:"access_token"`
	RefreshToken string `rose:"refresh_token"`
}

func (a *App) UserLogin(ctx *ctxx.Context, p UserLoginParams) (res UserLoginResponse, err error) {
	defer errs.DeferTrace(&err)()

	if p.WhatsAppNumber == "" && p.UserName == "" {
		return res, errs.ErrBadRequest.WithTrace().WithUserMessagef("Username or WhatsApp number is required")
	}

	where := []db.Where{}
	if p.WhatsAppNumber != "" {
		where = append(where, db.Where{
			Field: "whatsapp_number",
			Value: p.WhatsAppNumber,
		})
	} else if p.UserName != "" {
		where = append(where, db.Where{
			Field: "username",
			Value: p.UserName,
		})
	}

	users, _, err := a.d.UserUC.Get(ctx, nil, db.Params{
		Where: where,
		Pagination: &db.Pagination{
			Limit: 1,
		},
	})
	if err != nil {
		return res, err
	}
	if len(users) == 0 {
		return res, errs.ErrUnauthorized.WithTrace().WithUserMessagef("Invalid username/WhatsApp number or password")
	}
	user := users[0]

	ok, err := crypto.DefaultHasher.Compare(user.PasswordHash, user.PasswordSalt, []byte(p.Password))
	if err != nil {
		return res, err
	}
	if !ok {
		return res, errs.ErrUnauthorized.WithTrace().WithUserMessagef("Invalid username/WhatsApp number or password")
	}

	res.AccessToken, res.RefreshToken, err = a.makeTokenPair(user)
	if err != nil {
		return res, err
	}

	return res, nil
}

type RefreshTokenParams struct {
	RefreshToken string `rose:"refresh_token,rt,required="`
}

type RefreshTokenResponse struct {
	AccessToken string `rose:"jwt"`
}

func (a *App) RefreshToken(ctx *ctxx.Context, p RefreshTokenParams) (res RefreshTokenResponse, err error) {
	defer errs.DeferTrace(&err)()

	if p.RefreshToken == "" {
		return res, errs.ErrBadRequest.WithTrace().WithUserMessagef("Refresh token is required")
	}

	var claims entity.Claims
	_, err = a.d.JWT.Verify(p.RefreshToken, &claims)
	if err != nil {
		return res, err
	}
	if claims.TokenType != "refresh" {
		return res, errs.ErrBadRequest.WithTrace().WithUserMessagef("Invalid refresh token")
	}

	users, _, err := a.d.UserUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "id",
				Value: claims.UserID,
			},
		},
	})
	if err != nil {
		return res, err
	}
	if len(users) == 0 {
		return res, errs.ErrBadRequest.WithTrace().WithUserMessagef("Invalid refresh token")
	}
	user := users[0]

	res.AccessToken, err = a.makeAccessToken(user)
	if err != nil {
		return res, err
	}

	return res, nil
}
