package dbman

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elanticrypt0/dbman/console"
	"github.com/elanticrypt0/dbman/errors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type DBMan struct {
	rootPath         string
	configPath       string
	connection       map[string]DBConnection
	activeConnection []string
	Primary          *gorm.DB
	Secondary        *gorm.DB
	Security         *gorm.DB
}

func New() *DBMan {
	return &DBMan{
		rootPath:   "/api/",
		configPath: "/api/config",
		connection: make(map[string]DBConnection),
		Primary:    nil,
		Secondary:  nil,
		Security:   nil,
	}
}

// Load config from toml file
// This can load several database configurations
func (me *DBMan) SetRootPath(rootpath string) {
	me.rootPath = rootpath
}

// Load config from toml file
// This can load several database configurations
func (me *DBMan) LoadConfigToml(filepath string) {
	me.configPath = filepath
	configSlice := &DBConfigSlice{}
	LoadTomlFile(filepath, configSlice)

	// log.Fatalf("%+v", *configSlice)
	for _, config := range configSlice.Configurations {
		me.addConn(config)
	}
}

// Load database config from env file.
// For just one connection
func (me *DBMan) LoadConfigEnv() {
	envPath := "./.env"
	if ExitsFile(envPath) {
		err := godotenv.Load()
		if err != nil {
			log.Println(errors.FileNotLoaded(envPath))
		}
		connData := NewDBConfig(os.Getenv("DB_CONN_NAME"), os.Getenv("DB_ENGINE"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
		me.addConn(connData)
	} else {
		log.Println(errors.FileNotExistError(envPath))
	}

}

// Private function to add one connection to the connection slice
func (me *DBMan) addConn(connData DBConfig) {
	connNameLower := strings.ToLower(connData.ConnName)
	me.connection[connNameLower] = NewDBConn(connData)
}

func (me *DBMan) addActiveConn(conn string) {
	me.activeConnection = append(me.activeConnection, conn)
}

// Get one instance if the instance is connected
// otherwise returns nil and error.
// You need to call this from your code to use the gorm.DB
func (me *DBMan) GetInstance(name string) (*gorm.DB, error) {
	instance, err := me.getInstanceIfExists(name)
	if err != nil {
		return nil, err
	}
	// checks if the instance is connected
	if instance.IsConnected() {
		return instance.Instance, nil
	} else {
		return nil, instance.ErrConn
	}
}

// Private function
// Checks if the instance exist in connections slice
func (me *DBMan) getInstanceIfExists(name string) (*DBConnection, error) {
	name_lower := strings.ToLower(name)
	conn, exists := me.connection[name_lower]
	if exists {
		// if the instance exists checks if has no error
		return &conn, nil
	} else {
		errors.PrintStr(fmt.Sprintf("The connection %q (%q) does not exists \n", name, name_lower))
		return nil, errors.Instance("0", name, name_lower)
	}
}

// Checks if the DB connections has no errors
func (me *DBMan) IsDBOk(connName string) bool {
	conn := me.connection[connName]
	return conn.IsOk()
}

// Connects the selected connection
// If is the first connection
// by defaults is set as Primary
func (me *DBMan) Connect(name string) error {
	console.Print(fmt.Sprintf("Trying to connect to: %q", name))
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		errors.Print(err)
		return err
	}
	err = conn.Connect(me.rootPath)

	if err != nil {
		errors.Print(err)
		return err
	}

	me.connection[name] = *conn

	// automatic set of primary or secondary
	if me.Primary == nil {
		me.Primary = conn.Instance
	}
	if me.Secondary == nil {
		me.Secondary = conn.Instance
	}

	me.addActiveConn(conn.DBConfig.ConnName)
	console.Print(fmt.Sprintf("Connection stablishied to: %q", name))
	return nil
}

// Shortcut to the primary connection
func (me *DBMan) SetPrimary(name string) error {
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		return err
	}
	if conn.IsOk() {
		me.Primary = conn.Instance
		return nil
	} else {
		return conn.ErrConn
	}
}

// Shortcut to the secondary connection
func (me *DBMan) SetSecondary(name string) error {
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		return err
	}
	if conn.IsOk() {
		me.Secondary = conn.Instance
		return nil
	} else {
		return conn.ErrConn
	}
}

// Shortcut to the security or auth database
func (me *DBMan) SetSecurity(name string) error {
	conn, err := me.getInstanceIfExists(name)
	if err != nil {
		return err
	}
	if conn.IsOk() {
		me.Security = conn.Instance
		return nil
	} else {
		return conn.ErrConn
	}
}

/*  DEBUG Printing information  */

func (me *DBMan) PrintConnectionsList() {
	fmt.Println("")
	log.Println("DBMAN > Connections list")
	fmt.Println("---------------------------------------------")
	for _, connection := range me.connection {
		fmt.Printf("- %q \n", connection)
	}
	fmt.Printf("\n")

}

func (me *DBMan) PrintActiveConnectionsList() {
	fmt.Println("")
	log.Println("DBMAN > Active connections list")
	fmt.Println("----------------------------------------------------")
	for _, activeConn := range me.activeConnection {
		fmt.Printf("- %q \n", activeConn)
	}
	fmt.Printf("\n")

}

func (me *DBMan) CheckDefaultConnections() {

	fmt.Println("")
	log.Println("DBMAN > Check defaults connections list")
	fmt.Println("---------------------------------------------------------")

	fmt.Printf("- Primary (instance): %v \n", me.Primary)
	fmt.Printf("- Secondary (instance): %v \n", me.Secondary)
	fmt.Printf("- Security (instance): %v \n", me.Security)

}
