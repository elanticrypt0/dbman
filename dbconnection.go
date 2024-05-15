package dbman

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/k23dev/dbman/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnection struct {
	DBConfig DBConfig
	Instance *gorm.DB
	ErrConn  error
}

func NewDBConn(connData DBConfig) DBConnection {
	return DBConnection{
		DBConfig: connData,
		Instance: nil,
		ErrConn:  nil,
	}
}

func (me *DBConnection) IsOk() bool {
	return me.ErrConn == nil
}

func (me *DBConnection) IsConnected() bool {
	return me.Instance != nil
}

func (me *DBConnection) Connect() error {

	switch me.DBConfig.Engine {
	case "mysql":
		conn, err := me.connect2Mysql()
		if err != nil {
			return err
		}
		me.Instance = conn
	case "postgres":
		conn, err := me.connect2Postgres()
		if err != nil {
			return err
		}
		me.Instance = conn
	case "sqlite":
		conn, err := me.connect2SQLite()
		if err != nil {
			return err
		}
		me.Instance = conn
	default:
		log.Printf("The connections engine %q is not valid.\n", me.DBConfig.Engine)
		return nil
	}

	return nil
}

func (me *DBConnection) connect2Mysql() (*gorm.DB, error) {
	const dns = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dnsConfig := fmt.Sprintf(dns, me.DBConfig.User, me.DBConfig.Password, me.DBConfig.Host, me.DBConfig.PortAsStr, me.DBConfig.DBName)
	// connect to gorn
	conn, err := gorm.Open(mysql.Open(dnsConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		me.ErrConn = me.ReturnErrConn()
		return nil, me.ErrConn
	}
	return conn, nil
}

func (me *DBConnection) connect2Postgres() (*gorm.DB, error) {

	const dns = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable "
	dnsConfig := fmt.Sprintf(dns, me.DBConfig.Host, me.DBConfig.User, me.DBConfig.Password, me.DBConfig.DBName, me.DBConfig.PortAsStr)
	// connect to gorn
	conn, err := gorm.Open(postgres.Open(dnsConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		me.ErrConn = me.ReturnErrConn()
		return nil, me.ErrConn
	}
	// todo
	return conn, nil
}

func (me *DBConnection) connect2SQLite() (*gorm.DB, error) {
	dbname := me.DBConfig.DBName
	if ExitsFile(me.DBConfig.DBName) {
		// gorm create sqlite db
		conn, err := gorm.Open(sqlite.Open(dbname+".db"), &gorm.Config{})
		if err != nil {
			me.ErrConn = errors.Trying2ConnectSQLite("0", me.DBConfig.ConnName, me.DBConfig.Engine, me.DBConfig.DBName)
			return nil, me.ErrConn
		}
		return conn, nil
	} else {
		return nil, errors.Trying2ConnectSQLiteFileNotExists("0", me.DBConfig.DBName)
	}
}

func (me *DBConnection) ReturnErrConn() error {
	return errors.Trying2Connect("0", me.DBConfig.ConnName, me.DBConfig.Engine, me.DBConfig.Host, me.DBConfig.PortAsStr, me.DBConfig.User, me.DBConfig.Password, me.DBConfig.DBName)
}
