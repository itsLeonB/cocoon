package provider

import (
	"fmt"

	"github.com/itsLeonB/cocoon/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBs struct {
	dbConfig config.DB
	GormDB   *gorm.DB
}

func ProvideDBs(dbConfig config.DB) *DBs {
	dbs := &DBs{dbConfig, nil}
	dbs.openGormConnection()

	return dbs
}

func (d *DBs) Shutdown() error {
	db, err := d.GormDB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (d *DBs) getDSN() string {
	switch d.dbConfig.Driver {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.dbConfig.User,
			d.dbConfig.Password,
			d.dbConfig.Host,
			d.dbConfig.Port,
			d.dbConfig.Name,
		)
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s",
			d.dbConfig.Host,
			d.dbConfig.User,
			d.dbConfig.Password,
			d.dbConfig.Name,
			d.dbConfig.Port,
		)
	default:
		panic(fmt.Sprintf("unsupported SQLDB driver: %s", d.dbConfig.Driver))
	}
}

func (d *DBs) getGormDialector() gorm.Dialector {
	switch d.dbConfig.Driver {
	// case "mysql":
	// 	return mysql.Open(sqldb.getDSN())
	case "postgres":
		return postgres.Open(d.getDSN())
	default:
		panic(fmt.Sprintf("unsupported SQLDB driver: %s", d.dbConfig.Driver))
	}
}

func (d *DBs) openGormConnection() {
	db, err := gorm.Open(d.getGormDialector(), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("error opening GORM connection: %s", err.Error()))
	}

	d.GormDB = db
}
