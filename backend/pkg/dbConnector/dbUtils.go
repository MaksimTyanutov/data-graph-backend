package dbConnector

import (
	"data-graph-backend/pkg/properties"
	"database/sql"
	"errors"
	"fmt"
)

func (con *PSQLConnector) getMinEmployeeNum() (int, error) {
	var nullInt sql.NullInt32
	command := fmt.Sprintf("SELECT \"EmployeesNum\" FROM \"Company\" where \"EmployeesNum\" > 0 order by \"EmployeesNum\" LIMIT 1;")
	if err := con.db.QueryRow(command).Scan(&nullInt); err != nil {
		return 0, errors.New("Can't get min employee num. Error:" + err.Error())
	}
	return int(nullInt.Int32), nil
}

func (con *PSQLConnector) getMaxEmployeeNum() (int, error) {
	var nullInt sql.NullInt32
	command := fmt.Sprintf("SELECT \"EmployeesNum\" FROM \"Company\" where \"EmployeesNum\" > 0 order by \"EmployeesNum\" desc LIMIT 1;")
	if err := con.db.QueryRow(command).Scan(&nullInt); err != nil {
		return 0, errors.New("Can't get max employee num. Error:" + err.Error())
	}
	return int(nullInt.Int32), nil
}

func (con *PSQLConnector) getMinYearCompany() (string, error) {
	var nullString sql.NullString
	command := fmt.Sprintf("SELECT \"FoundationYear\" FROM \"Company\" where \"FoundationYear\" IS NOT NULL order by \"FoundationYear\" LIMIT 1;")
	if err := con.db.QueryRow(command).Scan(&nullString); err != nil {
		return "", errors.New("Can't get min year. Error:" + err.Error())
	}
	return nullString.String, nil
}

func (con *PSQLConnector) getMaxYearCompany() (string, error) {
	var nullString sql.NullString
	command := fmt.Sprintf("SELECT \"FoundationYear\" FROM \"Company\" where \"FoundationYear\" IS NOT NULL order by \"FoundationYear\" desc LIMIT 1;")
	if err := con.db.QueryRow(command).Scan(&nullString); err != nil {
		return "", errors.New("Can't get max year. Error:" + err.Error())
	}
	return nullString.String, nil
}

func (con *PSQLConnector) getMinDateProduct() (string, error) {
	var nullString sql.NullString
	command := fmt.Sprintf("SELECT \"Date\" FROM \"Projects\" where \"Date\" IS NOT NULL order by \"Date\" LIMIT 1")
	if err := con.db.QueryRow(command).Scan(&nullString); err != nil {
		return "", errors.New("Can't get min date for product. Error:" + err.Error())
	}
	return nullString.String, nil
}

func (con *PSQLConnector) getMaxDateProduct() (string, error) {
	var nullString sql.NullString
	command := fmt.Sprintf("SELECT \"Date\" FROM \"Projects\" where \"Date\" IS NOT NULL order by \"Date\" desc LIMIT 1")
	if err := con.db.QueryRow(command).Scan(&nullString); err != nil {
		return "", errors.New("Can't get max date for product. Error:" + err.Error())
	}
	return nullString.String, nil
}

func (con *PSQLConnector) GetMaxProductId() (int, error) {
	var nullInt sql.NullInt32
	command := fmt.Sprintf("select nodeid from getprojects() order by nodeid desc LIMIT 1")
	if err := con.db.QueryRow(command).Scan(&nullInt); err != nil {
		return 0, errors.New("Can't get max nodeId. Error:" + err.Error())
	}
	return int(nullInt.Int32), nil
}

func (con *PSQLConnector) SetIdShift() error {
	maxNodeId, err := con.GetMaxProductId()
	if err != nil {
		return errors.New("GetGraph(1): Can't get max nodeId: " + err.Error())
	}
	properties.CompanyIdShift = maxNodeId
	return nil
}
