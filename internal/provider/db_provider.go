package provider

import (
	"fmt"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/rotisserie/eris"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBs struct {
	dbConfig config.DB
	GormDB   *gorm.DB
}

func ProvideDBs(dbConfig config.DB) (*DBs, error) {
	dbs := &DBs{dbConfig, nil}
	if err := dbs.openGormConnection(); err != nil {
		return nil, err
	}
	return dbs, nil
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

func (d *DBs) openGormConnection() error {
	db, err := gorm.Open(d.getGormDialector(), &gorm.Config{})
	if err != nil {
		return eris.Wrap(err, "error opening gorm connection")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return eris.Wrap(err, "error returning sql DB")
	}

	sqlDB.SetMaxOpenConns(d.dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(d.dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(d.dbConfig.ConnMaxLifetime)

	d.GormDB = db
	return nil
}
