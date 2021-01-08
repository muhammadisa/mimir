package mimir

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConnectorBuilder struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func mySQLTestCall(db *sql.DB, testTable string) error {
	_, err := db.Exec(fmt.Sprintf("SELECT id FROM %s LIMIT 1", testTable))
	if err != nil {
		return err
	}
	return nil
}

func (dcb DBConnectorBuilder) MySQLConnect(testTable string) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=true",
		dcb.DBUser,
		dcb.DBPassword,
		dcb.DBHost,
		dcb.DBPort,
		dcb.DBName,
	)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	err = mySQLTestCall(db, testTable)
	if err != nil {
		return nil, err
	}
	return db, nil
}
