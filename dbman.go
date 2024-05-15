package dbman

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/k23dev/dbman/errors"
)

type DBMan struct {
	Connection map[string]DBConnection
	Primary    *DBConnection
	Secondary  *DBConnection
	Security   *DBConnection
}

func New() *DBMan {
	return &DBMan{
		Connection: make(map[string]DBConnection),
		Primary:    nil,
		Secondary:  nil,
		Security:   nil,
	}
}

func (me *DBMan) LoadConfigToml() {
	// todo
}

func (me *DBMan) LoadConfigEnv() {
	envPath := "./.env"
	if ExitsFile(envPath) {
		err := godotenv.Load()
		if err != nil {
			log.Println(errors.FileNotLoaded(envPath))
		}
		connData := NewDBConfig(os.Getenv("DB_CONN_NAME"), os.Getenv("DB_ENGINE"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		me.AddConn(connData)
	} else {
		log.Println(errors.FileNotExistError(envPath))
	}

}

func (me *DBMan) AddConn(connData DBConfig) {
	connNameLower := strings.ToLower(connData.ConnName)
	me.Connection[connNameLower] = NewDBConn(connData)
}

func (me *DBMan) GetInstance(name string) (*DBConnection, error) {
	instance, err := me.getInstanceIfExists(name)
	if err != nil {
		return nil, err
	}

	// checks if the instance is connected
	if instance.IsConnected() {
		return instance, nil
	} else {
		return nil, instance.ErrConn
	}
}

// internal funtion
func (me *DBMan) getInstanceIfExists(name string) (*DBConnection, error) {
	name_lower := strings.ToLower(name)
	conn, exists := me.Connection[name_lower]
	if exists {
		// if the instance exists checks if has no error
		return &conn, nil
	} else {
		log.Printf("The connection %s (%s) does not exists \n", name, name_lower)
		return nil, errors.Instance("0", name, name_lower)
	}
}

func (me *DBMan) IsDBOk(connName string) bool {
	conn := me.Connection[connName]
	return conn.IsOk()
}

func (me *DBMan) Connect(name string) error {
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		return err
	}
	err = conn.Connect()

	if err != nil {
		return err
	}

	return nil
}

func (me *DBMan) SetPrimary(name string) error {
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		return err
	}
	me.Primary = conn
	return nil
}

func (me *DBMan) SetSecondary(name string) error {
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		return err
	}
	me.Secondary = conn
	return nil
}

func (me *DBMan) SetSecurity(name string) error {
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		return err
	}
	me.Security = conn
	return nil
}
