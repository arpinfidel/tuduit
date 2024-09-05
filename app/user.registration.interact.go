package app

import (
	"fmt"
	"time"

	"github.com/arpinfidel/tuduit/entity"
	"github.com/arpinfidel/tuduit/pkg/crypto"
	"github.com/arpinfidel/tuduit/pkg/ctxx"
	"github.com/arpinfidel/tuduit/pkg/db"
	"github.com/arpinfidel/tuduit/pkg/errs"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

type OTPSendParams struct {
	WhatsAppNumber string `rose:"whatsapp_number,wa,required="`
}

type OTPSendResults struct {
}

func (a *App) OTPSend(ctx *ctxx.Context, p OTPSendParams) (res OTPSendResults, err error) {
	defer errs.DeferTrace(&err)()

	trx, err := a.d.DB.BeginTxx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer trx.Rollback()

	users, _, err := a.d.UserUC.Get(ctx, trx, db.Params{
		Where: []db.Where{
			{
				Field: "whatsapp_number",
				Value: p.WhatsAppNumber,
			},
		},
	})
	if err != nil {
		return res, err
	}
	if len(users) > 0 {
		return res, errs.ErrConflict.WithTrace().WithUserMessagef("WhatsApp number already exists")
	}

	otps, _, err := a.d.OTPUC.Get(ctx, trx, db.Params{
		Where: []db.Where{
			{
				Field: "whatsapp_number",
				Value: p.WhatsAppNumber,
			},
			{
				Field: "invalidated_at",
				Op:    db.GtOp,
				Value: time.Now(),
			},
		},
		Sort: []db.Sort{
			{
				Field: "invalidated_at",
				Asc:   false,
			},
		},
		Pagination: &db.Pagination{
			Limit: 3,
		},
	})
	if err != nil {
		return res, err
	}
	if len(otps) >= 3 { // too many requests
		return res, errs.ErrTooManyRequests.WithTrace().WithUserMessagef("Please try again in 15 minutes")
	}
	if len(otps) > 0 && time.Since(otps[0].CreatedAt) < 1*time.Minute { // rate limit
		return res, errs.ErrTooManyRequests.WithTrace().WithUserMessagef("Please try again in %s", time.Until(otps[0].CreatedAt.Add(1*time.Minute)).Round(time.Second))
	}

	otpCode := crypto.GenerateOTP(6)

	_, err = a.d.OTPUC.Create(ctx, trx, []entity.OTP{
		{
			WhatsAppNumber: p.WhatsAppNumber,
			OTP:            otpCode,
			InvalidatedAt:  time.Now().Add(15 * time.Minute),
			Token:          crypto.RandomString(32),
		},
	})
	if err != nil {
		return res, err
	}

	msg := fmt.Sprintf("Your OTP is %s. Please do not share the code with anyone.", otpCode)
	a.l.Debugf("Sending message: %s", msg)

	if err != nil {
		a.l.Errorf("Failed to connect to WhatsApp: %v", err)
		return res, err
	}
	_, err = a.d.WaClient.SendMessage(ctx, types.NewJID(p.WhatsAppNumber, types.DefaultUserServer), &waE2E.Message{
		Conversation: &msg,
	})
	if err != nil {
		return res, err
	}

	err = trx.Commit()
	if err != nil {
		return res, err
	}

	return res, nil
}

type OTPVerifyParams struct {
	WhatsAppNumber string `rose:"whatsapp_number,wa,required="`
	OTP            string `rose:"otp,required="`
}

type OTPVerifyResponse struct {
	Token string `rose:"token"`
}

func (a *App) OTPVerify(ctx *ctxx.Context, p OTPVerifyParams) (res OTPVerifyResponse, err error) {
	defer errs.DeferTrace(&err)()

	otps, _, err := a.d.OTPUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "whatsapp_number",
				Value: p.WhatsAppNumber,
			},
			{
				Field: "otp",
				Value: p.OTP,
			},
			{
				Field: "invalidated_at",
				Op:    db.GtOp,
				Value: time.Now(),
			},
		},
	})
	if err != nil {
		return res, err
	}
	if len(otps) == 0 {
		return res, errs.ErrBadRequest.WithTrace().WithUserMessagef("Invalid OTP")
	}

	res.Token = otps[0].Token

	return res, nil
}

type UserRegisterParams struct {
	Token          string `rose:"token,required="`
	WhatsAppNumber string `rose:"whatsapp_number,wa,required="`
	UserName       string `rose:"username,u,required="`
	Password       string `rose:"password,pw,required="`
}

type UserRegisterResponse struct {
	AccessToken  string `rose:"jwt"`
	RefreshToken string `rose:"refresh_token"`
}

func (a *App) UserRegister(ctx *ctxx.Context, p UserRegisterParams) (res UserRegisterResponse, err error) {
	otps, _, err := a.d.OTPUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "whatsapp_number",
				Value: p.WhatsAppNumber,
			},
			{
				Field: "token",
				Value: p.Token,
			},
			{
				Field: "invalidated_at",
				Op:    db.GtOp,
				Value: time.Now().Add(-15 * time.Minute),
			},
		},
	})
	if err != nil {
		return res, err
	}
	if len(otps) == 0 {
		return res, errs.ErrBadRequest.WithTrace().WithUserMessagef("Invalid token")
	}

	users, _, err := a.d.UserUC.Get(ctx, nil, db.Params{
		Where: []db.Where{
			{
				Field: "username",
				Value: p.UserName,
			},
		},
	})
	if err != nil {
		return res, err
	}
	if len(users) > 0 {
		return res, errs.ErrBadRequest.WithTrace().WithUserMessagef("Username already taken")
	}

	hashSalt, err := crypto.DefaultHasher.GenerateHash([]byte(p.Password), nil)
	if err != nil {
		return res, err
	}

	users, err = a.d.UserUC.Create(ctx, nil, []entity.User{
		{
			Username:       p.UserName,
			WhatsAppNumber: p.WhatsAppNumber,
			Name:           p.UserName,     // TODO: real name
			TimezoneStr:    "Asia/Jakarta", // TODO: user timezone
			PasswordHash:   hashSalt.Hash,
			PasswordSalt:   hashSalt.Salt,
		},
	})
	if err != nil {
		return res, err
	}
	user := users[0]

	res.AccessToken, res.RefreshToken, err = a.makeTokenPair(user)
	if err != nil {
		return res, err
	}

	return res, nil
}
