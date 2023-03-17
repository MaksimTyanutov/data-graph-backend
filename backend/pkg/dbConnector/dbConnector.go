package dbConnector

import (
	"data-graph-backend/pkg/dataStructers"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"strings"

	_ "github.com/lib/pq"

	"data-graph-backend/pkg/graphBuilder"
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
	logger *logrus.Logger
	db     *sql.DB
}

func (con *PSQLConnector) Ping() error {
	if err := con.db.Ping(); err != nil {
		return err
	}
	return nil
}

func NewConnection(config *properties.Config, logger *logrus.Logger) (*PSQLConnector, error) {
	db_, err := newDB(config)
	if err != nil {
		return nil, err
	}
	return &PSQLConnector{
		db:     db_,
		logger: logger,
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
			&c.ownernamessimilarity, &c.address, &c.iconpath, &c.posX, &c.posY); err != nil {
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
		var nullFloat sql.NullFloat64
		if err := rows.Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
			&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL, &nullFloat, &p.posX, &p.posY); err != nil {
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
		var nullFloat sql.NullFloat64
		if err := rows.Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
			&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL, &nullFloat, &p.posX, &p.posY); err != nil {
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
	if id < 1 {
		return nil, errors.New("Id company < 1: " + strconv.Itoa(id))
	}
	command := fmt.Sprintf("SELECT * From getcompanies(companyid := '%d')", id)
	c := new(Company)
	if err := con.db.QueryRow(command).Scan(&c.id, &c.name, &c.namesimilarity, &c.description, &c.descsimilarity,
		&c.employeeNum, &c.foundationyear, &c.companytypeenum, &c.companytypename, &c.ownerid, &c.ownername,
		&c.ownernamessimilarity, &c.address, &c.iconpath, &c.posX, &c.posY); err != nil {
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
		var nullFloat sql.NullFloat64
		if err := rows.Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
			&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL, &nullFloat, &p.posX, &p.posY); err != nil {
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
		Id:              company.Id + properties.CompanyIdShift,
		Name:            company.Name,
		Ceo:             company.OwnerName,
		Description:     company.Description,
		EmployeeNum:     company.EmployeeNum,
		FoundationYear:  company.FoundationYear,
		CompanyTypeName: company.CompanyTypeName,
		Products:        products,
		Svg:             company.IconPath,
	}
	return ci, nil
}

// Get info for product
func (con *PSQLConnector) GetProductInfo(id int) (*dataStructers.Product, error) {
	command := fmt.Sprintf("SELECT * From getprojects(searchnodeid := '%d')", id)
	p := new(Project)
	var nullFloat sql.NullFloat64
	if err := con.db.QueryRow(command).Scan(&p.nodeId, &p.projectId, &p.name, &p.nameSimilarity, &p.description, &p.version,
		&p.companyId, &p.projectTypesId, &p.projectTypesNames, &p.date, &p.url, &p.previousVersions, &p.pressURL, &nullFloat, &p.posX, &p.posY); err != nil {
		return nil, err
	}
	command = fmt.Sprintf("SELECT name From getcompanies(companyid := '%d')", int(p.companyId.Int32))
	var companyName sql.NullString
	if err := con.db.QueryRow(command).Scan(&companyName); err != nil {
		return nil, err
	}
	product, err := p.Transform()
	if err != nil {
		return nil, errors.New("graphBuilder:GetProjects(2). Can't transform project. " + err.Error())
	}
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
		return nil, errors.New("GetAllDepartments(2). Can't get departments from DB: " + err.Error())
	}
	return departments, nil
}

func (con *PSQLConnector) GetCompanyFilters() (*dataStructers.CompanyFilterPresets, error) {
	CF := dataStructers.CompanyFilterPresets{}

	departments, err := con.GetAllDepartments()
	if err != nil {
		return nil, errors.New("GetCompanyFilters(1): " + err.Error())
	}
	CF.Departments = departments
	companyNames, err := con.GetAllCompanyName()
	if err != nil {
		return nil, errors.New("GetCompanyFilters(2): " + err.Error())
	}
	CF.CompanyNames = companyNames
	ceoNames, err := con.GetAllCeoName()
	if err != nil {
		return nil, errors.New("GetCompanyFilters(3): " + err.Error())
	}
	CF.CeoNames = ceoNames

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

func (con *PSQLConnector) GetProductFilters() (*dataStructers.ProductFilterPresets, error) {
	productNames, err := con.GetAllProductName()
	if err != nil {
		return nil, err
	}
	minDate, err := con.getMinDateProduct()
	if err != nil {
		return nil, err
	}
	maxDate, err := con.getMaxDateProduct()
	if err != nil {
		return nil, err
	}
	PF := dataStructers.ProductFilterPresets{
		ProductNames: productNames,
		MinDate:      minDate,
		MaxDate:      maxDate,
	}
	return &PF, nil
}

func (con *PSQLConnector) GetFiltersIDCompany(companyFilter dataStructers.CompanyFilters) ([]int, error) {
	idArray := make([]int, 0)
	companyID := make([]int, 0)
	str := make([]string, 0)
	for _, el := range companyFilter.Departments {
		str = append(str, fmt.Sprintf("%d", el))
	}
	command := fmt.Sprintf("select id from getcompanies(namesearch := '%s', companytypeenums := '{%s}', ownersearch := '%s', begindate := '%s', enddate := '%s', employeescountbegin := '%d', employeescountend := '%d')",
		companyFilter.CompanyName, strings.Join(str, ", "), companyFilter.Ceo, companyFilter.MinDate, companyFilter.MaxDate, companyFilter.StartStaffSize, companyFilter.EndStaffSize)

	con.logger.Info("Executing sql query: " + command)
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("GetFiltersIDCompany(1). Can't get companies from DB: " + err.Error())
	}
	for rows.Next() {
		var id sql.NullInt32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		idArray = append(idArray, int(id.Int32)+properties.CompanyIdShift)
		companyID = append(companyID, int(id.Int32))
	}

	if len(companyID) == 0 {
		return idArray, nil
	}
	compId := make([]string, 0)
	for _, el := range companyID {
		compId = append(compId, fmt.Sprintf("%d", el))
	}
	con.logger.Info("Executing sql query: " + command)
	command = fmt.Sprintf("select nodeid from getprojects(companyids := '{%s}')", strings.Join(compId, ", "))
	rows2, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("GetFiltersIDCompany(2). Can't get projects from DB: " + err.Error())
	}
	for rows2.Next() {
		var id sql.NullInt32
		if err := rows2.Scan(&id); err != nil {
			return nil, err
		}
		idArray = append(idArray, int(id.Int32))
	}

	return idArray, nil
}

func (con *PSQLConnector) GetFiltersIDProduct(companyFilter dataStructers.ProductFilters) ([]int, error) {
	idArray := make([]int, 0)

	prjTypes := make([]string, 0)
	for _, el := range companyFilter.Departments {
		prjTypes = append(prjTypes, fmt.Sprintf("%d", el))
	}
	command := fmt.Sprintf("select nodeid from getprojects(namesearch := '%s', begindate := '%s', enddate := '%s',  searchprojecttypes := '{%s}', hasPressURL := '%t')",
		companyFilter.ProductName, companyFilter.MinDate, companyFilter.MaxDate, strings.Join(prjTypes, ", "), companyFilter.IsVerified)
	con.logger.Info("Executing sql query: " + command)
	rows2, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("GetFiltersID(2). Can't get projects from DB: " + err.Error())
	}
	for rows2.Next() {
		var id sql.NullInt32
		if err := rows2.Scan(&id); err != nil {
			return nil, err
		}
		idArray = append(idArray, int(id.Int32))
	}

	return idArray, nil
}

func (con *PSQLConnector) GetAllCompanyName() ([]string, error) {
	command := fmt.Sprintf("select distinct name from getcompanies()")
	companyNames := make([]string, 0)
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("GetAllCompanyName(1). Can't get execute command DB: " + err.Error())
	}
	for rows.Next() {
		var nullStr sql.NullString
		if err := rows.Scan(&nullStr); err != nil {
			return nil, err
		}
		companyNames = append(companyNames, nullStr.String)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("GetAllCompanyName(2). Can't get Names company from DB: " + err.Error())
	}
	return companyNames, nil
}

func (con *PSQLConnector) GetAllCeoName() ([]string, error) {
	command := fmt.Sprintf("select distinct ownername from getcompanies()")
	ceoNames := make([]string, 0)
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("GetAllCeoName(1). Can't get execute command DB: " + err.Error())
	}
	for rows.Next() {
		var nullStr sql.NullString
		if err := rows.Scan(&nullStr); err != nil {
			return nil, err
		}
		ceoNames = append(ceoNames, nullStr.String)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("GetAllCeoName(2). Can't get Names from DB: " + err.Error())
	}
	return ceoNames, nil
}

func (con *PSQLConnector) GetAllProductName() ([]string, error) {
	command := fmt.Sprintf("select distinct name from getprojects()")
	productNames := make([]string, 0)
	rows, err := con.db.Query(command)
	if err != nil {
		return nil, errors.New("GetAllProductName(1). Can't get execute command DB: " + err.Error())
	}
	for rows.Next() {
		var nullStr sql.NullString
		if err := rows.Scan(&nullStr); err != nil {
			return nil, err
		}
		productNames = append(productNames, nullStr.String)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("GetAllProductName(2). Can't get Names from DB: " + err.Error())
	}
	return productNames, nil
}

// GET ALL PROJECTS FOR GRAPH
func (con *PSQLConnector) GetProjectsGraph(minimized bool) ([]dataStructers.Project, error) {
	projects := make([]dataStructers.Project, 0)
	if !minimized {
		projectsDb, err := con.GetAllProjects()
		if err != nil {
			return nil, errors.New("graphBuilder:GetProjects(1). GetAllProjects don't work: " + err.Error())
		} else {
			for i := 0; i < len(projectsDb); i++ {
				project, err := projectsDb[i].Transform()
				if err != nil {
					return nil, errors.New("graphBuilder:GetProjects(2). Can't transform project. " + err.Error())
				}
				projects = append(projects, *project)
			}
		}
	} else {
		projectsDb, err := con.GetShortProjects()
		if err != nil {
			return nil, errors.New("graphBuilder:GetProjects(2). GetShortProjects don't work: " + err.Error())
		} else {
			for i := 0; i < len(projectsDb); i++ {
				project, err := projectsDb[i].Transform()
				if err != nil {
					return nil, errors.New("graphBuilder:GetProjects(2). Can't transform project. " + err.Error())
				}
				projects = append(projects, *project)
			}
		}
	}
	return projects, nil
}

// GET ALL COMPANIES  FOR GRAPH
func (con *PSQLConnector) GetCompaniesGraph() ([]dataStructers.Company, error) {
	companiesDb, err := con.GetAllCompanies()
	companies := make([]dataStructers.Company, 0)
	if err != nil {
		return nil, errors.New("graphBuilder:GetCompanies. GetAllCompanies don't work: " + err.Error())
	} else {
		for i := 0; i < len(companiesDb); i++ {
			company := companiesDb[i].Transform()
			companies = append(companies, company)
		}
	}
	return companies, nil
}

func (con *PSQLConnector) GetGraph(minimized bool) (*graphBuilder.Graph, error) {
	err := con.SetIdShift()
	if err != nil {
		return nil, err
	}
	companies, err := con.GetCompaniesGraph()
	if err != nil {
		return nil, errors.New("GetGraph(2): " + err.Error())
	}
	projects, err := con.GetProjectsGraph(minimized)
	if err != nil {
		return nil, errors.New("GetGraph(3): " + err.Error())
	}
	nodes := graphBuilder.TransformComp(companies)
	projectsTransformed := graphBuilder.TransformProj(projects)
	nodes = append(nodes, projectsTransformed...)
	links := graphBuilder.GetLinks(projects, minimized)
	return &graphBuilder.Graph{
		Nodes: nodes,
		Links: links,
	}, nil
}
