package vars

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/outcaste-io/ristretto"
	"gorm.io/gorm"
)

var (
	NetLogDatabase  bleve.Index
	AuthLogDatabase bleve.Index

	Database *gorm.DB
	// RedisCache *redis.Client

	AdminSession, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})

	BGImageCache, _ = ristretto.NewCache(&ristretto.Config{})
)
