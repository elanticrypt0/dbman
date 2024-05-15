package dbman

import (
	"log"
	"strconv"
)

type DBConfig struct {
	ConnName  string `toml:"connName"`
	Engine    string `toml:"engine"`
	Host      string `toml:"host"`
	Port      uint16 `toml:"port"`
	PortAsStr string
	User      string `toml:"user"`
	Password  string `toml:"password"`
	DBName    string `toml:"dbname"`
}

func NewDBConfig(connname, engine, host, port, user, passwd, dbname string) DBConfig {

	port_aux, err := strconv.Atoi(port)
	if err != nil {
		log.Println("Error converting DB_PORT as int")
	}

	return DBConfig{
		ConnName:  connname,
		Engine:    engine,
		Host:      host,
		Port:      uint16(port_aux),
		PortAsStr: port,
		User:      user,
		Password:  passwd,
		DBName:    dbname,
	}
}
