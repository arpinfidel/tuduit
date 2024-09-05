package entity

import "time"

type User struct {
	StdFields

	Name           string `json:"name"            db:"name"`
	Username       string `json:"username"        db:"username"`
	WhatsAppNumber string `json:"whatsapp_number" db:"whatsapp_number"`
	TimezoneStr    string `json:"timezone"        db:"timezone"`
	PasswordHash   []byte `json:"-"               db:"password_hash"`
	PasswordSalt   []byte `json:"-"               db:"password_salt"`
}

func (u *User) Timezone() *time.Location {
	timeZone, err := time.LoadLocation(u.TimezoneStr)
	if err != nil {
		panic(err)
	}

	return timeZone
}
