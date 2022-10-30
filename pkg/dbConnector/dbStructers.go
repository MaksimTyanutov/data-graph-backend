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
	id                  sql.NullInt32
	name                sql.NullString
	description         sql.NullString
	version             sql.NullString
	companyId           sql.NullInt32
	projectId           sql.NullInt32
	projectVersionIndex sql.NullInt32
	date                sql.NullString
	lastNodeIds         []uint8
	hasTwoInputs        sql.NullBool
	projectTypeIds      []uint8
	url                 sql.NullString
}

func (p *Project) GetName() string {
	return p.name.String
}
