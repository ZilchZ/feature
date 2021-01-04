package dao

import (
	"fmt"
	"github.com/feature/conf"
	"github.com/feature/sdk/log"
	"github.com/go-xorm/xorm"
)

var DBX *xorm.Engine

func Init() {
	if DBX == nil {
		config := conf.Config.Database
		connStr := fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
			config.Host, config.Port, config.User, config.Password, config.DBName)
		DBX, err := xorm.NewEngine(config.Driver, connStr)
		if err != nil {
			log.Logger.Error(err)

		} else {
			DBX.ShowSQL(conf.Config.Log.DevModel)
		}

	}
}
