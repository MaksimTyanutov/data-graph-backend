package dataStructers

type Company struct {
	id             int
	name           string
	description    string
	employeeNum    int
	foundationYear string
	companyTypeIds []int
	owner          string
	address        string
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
