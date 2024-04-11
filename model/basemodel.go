package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string ` gorm:"type:varchar(36);primarykey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (o *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}

	// o.CreatedAt = time.Now().In(vars.TZ)

	return
}

// func (o *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
// 	o.UpdatedAt = time.Now().In(vars.TZ)
// 	return
// }

// func (o *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
// 	o.DeletedAt = gorm.DeletedAt{
// 		Time:  time.Now().In(vars.TZ),
// 		Valid: true,
// 	}
// 	return
// }

type CompositeBaseModels struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type JsonType []byte

func (v *JsonType) GetMap() any {
	var data any
	json.Unmarshal(*v, &data)
	return data
}
