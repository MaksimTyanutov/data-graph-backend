package dbConnector

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Company struct {
	id                   sql.NullInt32
	name                 sql.NullString
	namesimilarity       sql.NullInt32
	description          sql.NullString
	descsimilarity       sql.NullInt32
	employeeNum          sql.NullInt32
	foundationyear       sql.NullString
	companytypeenum      []uint8
	companytypename      []uint8
	ownerid              sql.NullInt32
	ownername            sql.NullString
	ownernamessimilarity sql.NullInt32
	address              sql.NullString
	iconpath             sql.NullString
}

func (c *Company) GetName() string {
	return c.name.String
}

type Project struct {
	id                  int
	projectId           int
	projectTypeIds      []int
	name                string
	companyId           int
	version             string
	projectVersionIndex int
	date                string
	lastNodeIds         []int
	hasTwoInputs        bool
	url                 string
}
