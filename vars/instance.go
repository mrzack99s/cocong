package vars

import (
	"github.com/outcaste-io/ristretto"
	"gorm.io/gorm"
)

var (
	Database         *gorm.DB
	InMemoryDatabase *gorm.DB

	AdminSession, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})

	BGImageCache, _ = ristretto.NewCache(&ristretto.Config{})
)
