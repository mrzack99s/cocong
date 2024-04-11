package model

type Directory struct {
	BaseModel

	Name          string
	Enable        bool
	MaxConcurrent int64

	BandwidthID *string    `gorm:"type:varchar(36)"`
	Bandwidth   *Bandwidth `gorm:"references:id"`

	Users []User `json:"omitempty" valid:"-" gorm:"foreignKey:DirectoryID"`
}
