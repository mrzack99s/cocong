package model

type Administrator struct {
	BaseModel

	Name     string
	Username string
	Hashed   string
	Enable   bool
}
