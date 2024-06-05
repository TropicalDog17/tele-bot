package types

type Settings struct {
	Language        string `json:"language" db:"language"`
	Ticker          string `json:"ticker" db:"ticker"`
	PasswordEnabled bool   `json:"password_enabled" db:"password_enabled"`
}

func NewSettings() *Settings {
	return &Settings{
		Language:        "en",
		Ticker:          "ATOMINJ",
		PasswordEnabled: true,
	}
}
