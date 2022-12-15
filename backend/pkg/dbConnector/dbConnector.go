package dbConnector

import (
	"data-graph-backend/pkg/dataStructers"
	"database/sql"
	"errors"
	"fmt"
	"log"

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

func (con *PSQLConnector) Test() (string, error) {
	var str string
	command := fmt.Sprintf("SELECT name From getcompanies(companyid=>'%d')", 1)
	if err := con.db.QueryRow(command).Scan(&str); err != nil {
		return "", err
	}
	return str, nil
}

func (con *PSQLConnector) GetNumberCompanies() (int, error) {
	var total int
	command := fmt.Sprintf("SELECT COUNT(*) From \"Company\"")
	if err := con.db.QueryRow(command).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (con *PSQLConnector) GetAllCompanies() ([]Company, error) {
	companies := make([]Company, 0)
	command := fmt.Sprintf("SELECT * From getcompanies()")
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("Can't execute command: " + command + "; " + err.Error())
	}
	for rows.Next() {
		c := new(Company)
		if err := rows.Scan(&c.id, &c.name, &c.namesimilarity, &c.description, &c.descsimilarity,
			&c.employeeNum, &c.foundationyear, &c.companytypeenum, &c.companytypename, &c.ownerid, &c.ownername,
			&c.ownernamessimilarity, &c.address, &c.iconpath); err != nil {
			return nil, errors.New("Can't read company info: " + err.Error())
		}
		companies = append(companies, *c)
	}
	return companies, nil
}

func (con *PSQLConnector) GetNumberProjects() (int, error) {
	var total int
	command := fmt.Sprintf("SELECT COUNT(*) From \"Projects\"")
	if err := con.db.QueryRow(command).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (con *PSQLConnector) GetAllProjects() ([]Project, error) {
	projects := make([]Project, 0)
	command := fmt.Sprintf("SELECT * FROM getprojects(namesearch => '')")
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := new(Project)
		if err := rows.Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
			&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL); err != nil {
			return nil, err
		}
		projects = append(projects, *p)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err.Error())
	}
	return projects, nil
}

func (con *PSQLConnector) GetShortProjects() ([]Project, error) {
	projects := make([]Project, 0)
	command := fmt.Sprintf("SELECT * FROM getprojects(shortform := true)")
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := new(Project)
		if err := rows.Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
			&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL); err != nil {
			return nil, err
		}
		projects = append(projects, *p)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err.Error())
	}
	return projects, nil
}

// Get info for company
func (con *PSQLConnector) GetCompanyInfo(id int) (*dataStructers.CompanyInfo, error) {
	id = id - properties.CompanyIdShift
	command := fmt.Sprintf("SELECT * From getcompanies(companyid := '%d')", id)
	c := new(Company)
	if err := con.db.QueryRow(command).Scan(&c.id, &c.name, &c.namesimilarity, &c.description, &c.descsimilarity,
		&c.employeeNum, &c.foundationyear, &c.companytypeenum, &c.companytypename, &c.ownerid, &c.ownername,
		&c.ownernamessimilarity, &c.address, &c.iconpath); err != nil {
		return nil, errors.New("Can't execute command: " + command + ". Error:" + err.Error())
	}
	company := c.Transform()
	command = fmt.Sprintf("SELECT * From getprojects(companyids := '{%d}')", id)
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, err
	}
	products := make([]dataStructers.ProductShort, 0)
	for rows.Next() {
		p := new(Project)
		if err := rows.Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
			&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL); err != nil {
			return nil, err
		}
		prod := dataStructers.ProductShort{
			Id:         int(p.nodeId.Int32),
			Name:       p.name.String,
			Year:       p.date.String,
			IsVerified: p.pressURL.String != "",
		}
		products = append(products, prod)
	}
	ci := &dataStructers.CompanyInfo{
		Id:              company.Id,
		Name:            company.Name,
		Ceo:             company.OwnerName,
		Description:     company.Description,
		EmployeeNum:     company.EmployeeNum,
		FoundationYear:  company.FoundationYear,
		CompanyTypeName: company.CompanyTypeName,
		Products:        products,
	}
	return ci, nil
}

// Get info for product
func (con *PSQLConnector) GetProductInfo(id int) (*dataStructers.Product, error) {
	command := fmt.Sprintf("SELECT * From getprojects(searchnodeid := '%d')", id)
	p := new(Project)
	if err := con.db.QueryRow(command).Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
		&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL); err != nil {
		return nil, err
	}
	command = fmt.Sprintf("SELECT name From getcompanies(companyid := '%d')", int(p.companyId.Int32))
	var companyName sql.NullString
	if err := con.db.QueryRow(command).Scan(&companyName); err != nil {
		return nil, err
	}
	product := p.Transform()
	departments := make([]dataStructers.Department, 0)
	for i := 0; i < len(product.ProjectTypes); i++ {
		dep := dataStructers.Department{
			Id:   int(rune(int(p.projectTypesId[i*2]))),
			Name: product.ProjectTypes[i],
		}
		departments = append(departments, dep)
	}
	pi := &dataStructers.Product{
		Id:      product.Id,
		Name:    product.Name,
		Version: product.Version,
		Company: struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}{
			Id:   product.CompanyId,
			Name: companyName.String,
		},
		Link:        product.PressURL,
		Description: product.Description,
		Svg:         product.Url,
		Year:        product.Date,
		Departments: departments,
	}
	return pi, nil
}

func (con *PSQLConnector) GetAllDepartments() ([]dataStructers.Department, error) {
	command := fmt.Sprintf("select * from \"Departments\"")
	departments := make([]dataStructers.Department, 0)
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("GetAllDepartments(1). Can't get departments from DB: " + err.Error())
	}
	for rows.Next() {
		d := new(Department)
		if err := rows.Scan(&d.id, &d.name); err != nil {
			return nil, err
		}
		department := d.Transform()
		departments = append(departments, department)
	}
	if err := rows.Err(); err != nil {
		log.Println("GetAllDepartments(2). Can't get departments from DB: " + err.Error())
	}
	return departments, nil
}

func (con *PSQLConnector) GetCompanyFilters() (*dataStructers.CompanyFilters, error) {
	CF := dataStructers.CompanyFilters{}

	departments, err := con.GetAllDepartments()
	if err != nil {
		return nil, errors.New("GetCompanyFilters(3): " + err.Error())
	}
	CF.Departments = departments

	CF.MinStaffSize, err = con.getMinEmployeeNum()
	if err != nil {
		return nil, err
	}
	CF.MaxStaffSize, err = con.getMaxEmployeeNum()
	if err != nil {
		return nil, err
	}
	CF.MinDate, err = con.getMinYearCompany()
	if err != nil {
		return nil, err
	}
	CF.MaxDate, err = con.getMaxYearCompany()
	if err != nil {
		return nil, err
	}

	return &CF, nil
}

func (con *PSQLConnector) GetProductFilters() (*dataStructers.ProductFilters, error) {
	minDate, err := con.getMinDateProduct()
	if err != nil {
		return nil, err
	}
	maxDate, err := con.getMaxDateProduct()
	if err != nil {
		return nil, err
	}
	PF := dataStructers.ProductFilters{
		MinDate: minDate,
		MaxDate: maxDate,
	}
	return &PF, nil
}
