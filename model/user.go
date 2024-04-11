package model

type User struct {
	BaseModel

	Name             string
	Enable           bool
	UserID           string
	Username         string
	Hashed           string
	FailedLoginCount int64

	DirectoryID *string    `gorm:"type:varchar(36)"`
	Directory   *Directory `gorm:"references:id"`
}
