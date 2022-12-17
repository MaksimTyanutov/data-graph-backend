package dataStructers

import (
	"encoding/json"
	"log"
)

type Company struct {
	Id              int
	Name            string
	Description     string
	EmployeeNum     int
	FoundationYear  string
	CompanyTypeName []string
	OwnerName       string
	Address         string
	IconPath        string
}

func (c *Company) JSON() string {
	marshal, err := json.Marshal(c)
	if err != nil {
		log.Fatal("Error while marshalling: " + err.Error())
	}
	return string(marshal)
}

func (c *Company) SetId(id int) {
	c.Id = id
}

func (c *Company) SetName(name string) {
	c.Name = name
}

func (c *Company) SetDescription(description string) {
	c.Description = description
}

func (c *Company) SetEmployeeNum(employeeNum int) {
	c.EmployeeNum = employeeNum
}

func (c *Company) SetFoundationYear(foundationYear string) {
	c.FoundationYear = foundationYear
}

func (c *Company) SetCompanyTypeName(companyTypeName []string) {
	c.CompanyTypeName = companyTypeName
}

func (c *Company) SetOwnerName(ownerName string) {
	c.OwnerName = ownerName
}

func (c *Company) SetAddress(address string) {
	c.Address = address
}

func (c *Company) SetIconPath(iconPath string) {
	c.IconPath = iconPath
}

type Project struct {
	Id              int
	ProjectId       int
	Name            string
	Description     string
	Version         string
	CompanyId       int
	ProjectTypes    []string
	Date            string
	Url             string
	PreviousNodeIds []int
	PressURL        string
}

func (p *Project) JSON() string {
	marshal, err := json.Marshal(p)
	if err != nil {
		log.Fatal("Error while marshalling: " + err.Error())
	}
	return string(marshal)
}

type Product struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Company struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"company"`
	Link        string       `json:"link"`
	Description string       `json:"description"`
	Svg         string       `json:"svg"`
	Year        string       `json:"year"`
	Departments []Department `json:"departments"`
}

type CompanyInfo struct {
	Id              int            `json:"id"`
	Name            string         `json:"name"`
	Ceo             string         `json:"ceo"`
	Description     string         `json:"description"`
	EmployeeNum     int            `json:"staffSize"`
	FoundationYear  string         `json:"year"`
	CompanyTypeName []string       `json:"departments"`
	Products        []ProductShort `json:"products"`
}

type ProductShort struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Year       string `json:"year"`
	IsVerified bool   `json:"isVerified"`
}

type Department struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TimelineCompany struct {
	Company struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Year string `json:"year"`
	} `json:"company"`
	Product struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Year string `json:"year"`
	} `json:"product"`
}

type TimelineProducts struct {
	Product1 struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Year string `json:"year"`
	} `json:"product-1"`
	Product2 struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Year string `json:"year"`
	} `json:"product-2"`
}

type FilterPresets struct {
	CompanyFilters CompanyFilterPresets `json:"companyFilters"`
	ProductFilters ProductFilterPresets `json:"productFilters"`
}

type CompanyFilterPresets struct {
	CompanyNames []string     `json:"companyNames"`
	CeoNames     []string     `json:"ceoNames"`
	MinStaffSize int          `json:"minStaffSize"`
	MaxStaffSize int          `json:"maxStaffSize"`
	MinDate      string       `json:"minDate"`
	MaxDate      string       `json:"maxDate"`
	Departments  []Department `json:"departments"`
}

type ProductFilterPresets struct {
	ProductNames []string `json:"productNames"`
	MinDate      string   `json:"minDate"`
	MaxDate      string   `json:"maxDate"`
}

type Filters struct {
	CompanyFilters CompanyFilters `json:"companyFilters"`
	ProductFilters ProductFilters `json:"productFilters"`
}

type CompanyFilters struct {
	CompanyName    string `json:"companyName"`
	Departments    []int  `json:"departments"`
	Ceo            string `json:"ceo"`
	MinDate        string `json:"minDate"`
	MaxDate        string `json:"maxDate"`
	StartStaffSize int    `json:"startStaffSize"`
	EndStaffSize   int    `json:"endStaffSize"`
}

type ProductFilters struct {
	ProductName string `json:"productName"`
	MinDate     string `json:"minDate"`
	MaxDate     string `json:"maxDate"`
	Departments []int  `json:"departments"`
	IsVerified  bool   `json:"isVerified"`
}
