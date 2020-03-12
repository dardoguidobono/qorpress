package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/qorpress/qorpress/core/l10n"
	"github.com/qorpress/qorpress/core/media"
	"github.com/qorpress/qorpress/core/publish2"
	"github.com/qorpress/qorpress/core/sorting"
	"github.com/qorpress/qorpress/core/validations"
	"github.com/qorpress/qorpress/pkg/config"
)

// DB Global DB connection
var DB *gorm.DB

func init() {
	var err error

	dbConfig := config.Config.DB

	pp.Println("dbConfig:", dbConfig)

	if config.Config.DB.Adapter == "mysql" {
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
		DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		// DB = DB.Set("gorm:table_options", "CHARSET=utf8")
	} else if config.Config.DB.Adapter == "postgres" {
		DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name))
	} else if config.Config.DB.Adapter == "sqlite" {
		DB, err = gorm.Open("sqlite3", fmt.Sprintf("%v/%v", os.TempDir(), dbConfig.Name))
	} else {
		panic(errors.New("not supported database adapter"))
	}

	if err == nil {
		if os.Getenv("DEBUG") != "" {
			DB.LogMode(true)
		}

		l10n.RegisterCallbacks(DB)
		sorting.RegisterCallbacks(DB)
		validations.RegisterCallbacks(DB)
		media.RegisterCallbacks(DB)
		publish2.RegisterCallbacks(DB)
	} else {
		panic(err)
	}
}
