package main

import (
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/properties"
	"fmt"
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

	companies, err := dbConnection.GetAllCompanies()
	if err != nil {
		log.Print("GetAllCompanies don't work: ", err.Error())
	} else {
		for i := 0; i < len(companies); i++ {
			fmt.Println(companies[i].GetName())
		}
	}

	projects, err := dbConnection.GetAllProjects()
	if err != nil {
		log.Print("GetAllProjects don't work: ", err.Error())
	} else {
		for i := 0; i < len(projects); i++ {
			fmt.Println(projects[i].GetName())
		}
	}
}
