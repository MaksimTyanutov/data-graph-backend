package dbConnector

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"data-graph-backend/pkg/properties"
)

func newDB(config *properties.Config) (*sql.DB, error) {
	dbHost := config.DbSettings.DbHost
	dbName := config.DbSettings.DbName
	dbUsername := config.DbSettings.DbUsername
	dbPassword := config.DbSettings.DbPassword
	dbPort := config.DbSettings.DbPort

	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=disable", dbHost, dbName, dbUsername, dbPassword,
		dbPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type PSQLConnector struct {
	db *sql.DB
}

func NewConnection(config *properties.Config) (*PSQLConnector, error) {
	db_, err := newDB(config)
	if err != nil {
		return nil, err
	}
	return &PSQLConnector{
		db: db_,
	}, nil
}

func (c *PSQLConnector) Test() (string, error) {
	var str string
	command := fmt.Sprintf("SELECT name From getcompanies(companyid=>'%d')", 1)
	if err := c.db.QueryRow(command).Scan(&str); err != nil {
		return "", err
	}
	return str, nil
}

func (c *PSQLConnector) GetNumberCompanies() (int, error) {
	var total int
	command := fmt.Sprintf("SELECT COUNT(*) From \"Company\"")
	if err := c.db.QueryRow(command).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
