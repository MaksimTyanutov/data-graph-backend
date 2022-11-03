package main

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/properties"
	"data-graph-backend/pkg/utils"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Incorrect number of arguments: ", len(os.Args), ". Intended number - 2")
	}
	config := properties.GetConfig(os.Args[1])
	dbConnection, err := dbConnector.NewConnection(config)
	if err != nil {
		log.Fatal("Can't connect to db - ", err.Error())
	}

	//GET ALL COMPANIES
	jsonStr := "["
	companiesDb, err := dbConnection.GetAllCompanies()
	companies := make([]dataStructers.Company, 0)
	if err != nil {
		log.Print("GetAllCompanies don't work: ", err.Error())
	} else {
		for i := 0; i < len(companiesDb); i++ {
			company := companiesDb[i].Transform()
			companies = append(companies, company)
			jsonStr = jsonStr + company.JSON() + ","
		}
	}
	jsonStr = jsonStr + "]"
	utils.ToFile(jsonStr, "Companies")

	//GET ALL PROJECTS
	jsonStr = "["
	projectsDb, err := dbConnection.GetAllProjects()
	if err != nil {
		log.Print("GetAllProjects don't work: ", err.Error())
	} else {
		projects := make([]dataStructers.Project, 0)
		for i := 0; i < len(projectsDb); i++ {
			project := projectsDb[i].Transform()
			projects = append(projects, project)
			jsonStr = jsonStr + project.JSON() + ","
		}
	}
	jsonStr = jsonStr + "]"
	utils.ToFile(jsonStr, "Projects")
}
