package setup

import (
	"fmt"
	"time"

	"github.com/mrzack99s/cocong/constants"
	"github.com/mrzack99s/cocong/drivers/sqlite"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Database() {
	var err error

	dbname := "cocong.db"
	if !vars.SYS_DEBUG {
		dbname = constants.APP_DIR + "/cocong.db"
	}

	if !vars.SYS_DEBUG {
		key := "Brey@crazu!risweqI==4!boyu?oqE*As+Spenu7i#OtaYif!i$RokE=h*PlDr@p"
		dbnameWithDSN := dbname + fmt.Sprintf("?_pragma_key=%s&_pragma_cipher_page_size=4096", key)
		vars.Database, err = gorm.Open(sqlite.Open(dbnameWithDSN), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().In(vars.TZ)
			},
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			panic(err)
		}
	} else {
		vars.Database, err = gorm.Open(sqlite.Open(dbname), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			panic(err)
		}
	}

	inmemoryDBName := "file::memory:?cache=shared"
	if vars.SYS_DEBUG {
		inmemoryDBName = "inmemory.db"
	}
	vars.InMemoryDatabase, err = gorm.Open(sqlite.Open(inmemoryDBName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	dbmigrate()
	dbinit()
}

func dbmigrate() {
	models := []any{
		model.Bandwidth{},
		model.LoginLog{},
		model.LogoutLog{},
		model.Administrator{},
		model.AdministratorLoginLog{},
		model.Directory{},
		model.User{},
		model.NetworkLog{},
	}

	inmemoryModels := []any{
		inmemory_model.Session{},
	}

	vars.Database.AutoMigrate(models...)
	vars.InMemoryDatabase.AutoMigrate(inmemoryModels...)
}

func dbinit() {
	if err := vars.Database.Select("id").First(&model.Administrator{}).Error; err != nil {
		vars.Database.Create(&model.Administrator{
			Name:     "Initial Administrator",
			Enable:   true,
			Username: "administrator",
			Hashed:   utils.Sha512encode("P@ssw0rd"),
		})
	}

	// Pre-defined bandwidth

	if err := vars.Database.Select("id").First(&model.Bandwidth{}).Error; err != nil {

		name := "1Gbps/1Gbps"
		vars.Database.Create(&model.Bandwidth{
			Name:          name,
			DownloadLimit: 1000,
			UploadLimit:   1000,
		})

		name = "1Gbps/500Mbps"
		vars.Database.Create(&model.Bandwidth{
			Name:          name,
			DownloadLimit: 1000,
			UploadLimit:   500,
		})

		name = "500Mbps/500Mbps"
		vars.Database.Create(&model.Bandwidth{
			Name:          name,
			DownloadLimit: 500,
			UploadLimit:   500,
		})
	}
}
