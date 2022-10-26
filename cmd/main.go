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
	str, err := dbConnection.Test()
	if err != nil {
		log.Print("Test don't work: ", err.Error())
	} else {
		fmt.Print(str)
	}
}
