package entity

import "time"

type User struct {
	StdFields

	Name           string         `json:"name"            db:"name"`
	Username       string         `json:"username"        db:"username"`
	WhatsappNumber string         `json:"whatsapp_number" db:"whatsapp_number"`
	TimeZoneStr    string         `json:"timezone"        db:"timezone"`
	TimeZone       *time.Location `json:"-"               db:"-"`
}

func (u *User) Parse() error {
	timeZone, err := time.LoadLocation(u.TimeZoneStr)
	if err != nil {
		return err
	}
	u.TimeZone = timeZone

	return nil
}
