package gaarx

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type (
	db struct {
		database        *gorm.DB
		migrateEntities []interface{}
	}
)

func (d *db) GetDB() *gorm.DB {
	return d.database
}

func (d *db) MigrateEntities(entities ...interface{}) {
	for _, e := range entities {
		d.database.AutoMigrate(e)
	}
}

func GetConnString(user, pass, host, port, dbName string) string {
	return user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4,utf8&parseTime=true&sql_mode=ansi"
}
