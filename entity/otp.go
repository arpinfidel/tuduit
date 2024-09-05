package entity

import "time"

type OTP struct {
	StdFields

	WhatsAppNumber string    `db:"whatsapp_number" json:"whatsapp_number"`
	OTP            string    `db:"otp"             json:"otp"`
	Token          string    `db:"token"           json:"token"`
	InvalidatedAt  time.Time `db:"invalidated_at"  json:"invalidated_at"`
}
