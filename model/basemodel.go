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

type BaseModelWithoutAutoID struct {
	ID        string ` gorm:"type:varchar(36);primarykey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RedisBaseModel struct {
	ID string ` gorm:"type:varchar(36);primarykey;"`
}

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
