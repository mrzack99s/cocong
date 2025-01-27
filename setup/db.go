package setup

import (
	"fmt"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/mrzack99s/cocong/constants"
	"github.com/mrzack99s/cocong/drivers/sqlite"
	"github.com/mrzack99s/cocong/model"
	"github.com/mrzack99s/cocong/utils"
	"github.com/mrzack99s/cocong/vars"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Database() {
	var err error

	dbname := "cocong.db"
	netlog_dbname := "cocong-netlog.dblog"
	authlog_dbname := "cocong-authlog.dblog"
	if !vars.SYS_DEBUG {
		if vars.Config.CustomDBDir != "" {
			dbname = vars.Config.CustomDBDir + "/cocong.db"
			netlog_dbname = vars.Config.CustomDBDir + "/" + netlog_dbname
			authlog_dbname = vars.Config.CustomDBDir + "/" + authlog_dbname
		} else {
			dbname = constants.APP_DIR + "/cocong.db"
			netlog_dbname = constants.APP_DIR + "/" + netlog_dbname
			authlog_dbname = constants.APP_DIR + "/" + authlog_dbname
		}

	}

	dsnTail := "?cache=shared&_busy_timeout=5000&_journal_mode=WAL&_synchronous=OFF&_temp_store=MEMORY&_cache_size=-1000000"
	if vars.Config.DBCacheSize > 0 {
		dsnTail = fmt.Sprintf("?cache=shared&_busy_timeout=5000&_journal_mode=WAL&_synchronous=OFF&_temp_store=MEMORY&_cache_size=-%d", vars.Config.DBCacheSize)
	}

	if !vars.SYS_DEBUG {

		dbnameWithDSN := dbname + dsnTail
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

		if _, err := os.Stat(netlog_dbname); os.IsNotExist(err) {
			vars.NetLogDatabase, err = bleve.NewUsing(netlog_dbname, bleve.NewIndexMapping(), "scorch", "scorch", nil)
			if err != nil {
				panic(err)
			}

		} else {

			vars.NetLogDatabase, err = bleve.Open(netlog_dbname)
			if err != nil {
				panic(err)
			}

		}

		if _, err := os.Stat(authlog_dbname); os.IsNotExist(err) {
			vars.AuthLogDatabase, err = bleve.NewUsing(authlog_dbname, bleve.NewIndexMapping(), "scorch", "scorch", nil)
			if err != nil {
				panic(err)
			}

		} else {

			vars.AuthLogDatabase, err = bleve.Open(authlog_dbname)
			if err != nil {
				panic(err)
			}

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

	sqlDB, err := vars.Database.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(0) // Maximum number of open connections
	sqlDB.SetMaxIdleConns(0) // Maximum number of idle connections

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
	}

	vars.Database.AutoMigrate(models...)
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
