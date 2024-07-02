package entity

type User struct {
	StdFields

	Name           string `json:"name"            db:"name"`
	Username       string `json:"username"        db:"username"`
	WhatsappNumber string `json:"whatsapp_number" db:"whatsapp_number"`
}
