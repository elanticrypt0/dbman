package errors

import (
	"errors"
	"fmt"
)

type DBManErr struct {
	StatusCode string
	Err        error
}

func (me *DBManErr) Error() string {
	return fmt.Sprintf("status %s: err %v", me.StatusCode, me.Err)
}

func Generic(code, msg string) *DBManErr {
	return &DBManErr{
		StatusCode: code,
		Err:        errors.New(msg),
	}
}

func Connection(code, conn string) error {
	return &DBManErr{
		StatusCode: code,
		Err:        errors.New(GetConnectionError(conn)),
	}
}

func ConnectionFails(code string) error {
	return &DBManErr{
		StatusCode: code,
		Err:        errors.New(GetFailedToConnect()),
	}
}

func Trying2Connect(code, connname, engine, host, port, user, password, dbname string) error {
	return &DBManErr{
		StatusCode: code,
		Err:        errors.New(GetTrying2ConnectError(connname, engine, host, port, user, password, dbname)),
	}
}

func Trying2ConnectSQLite(code, connname, engine, dbname string) error {
	return &DBManErr{
		StatusCode: code,
		Err:        errors.New(GetTrying2ConnectSQLiteError(connname, engine, dbname)),
	}
}

func Trying2ConnectSQLiteFileNotExists(code, filepath string) error {
	return &DBManErr{
		StatusCode: code,
		Err:        errors.New(FileNotExistError(filepath)),
	}
}

func Instance(code, name, name_lower string) error {
	return &DBManErr{
		StatusCode: code,
		Err:        errors.New(GetInstanceError(name, name_lower)),
	}
}
