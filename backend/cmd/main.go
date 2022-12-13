package main

import (
	"data-graph-backend/pkg/apiServer"
	"data-graph-backend/pkg/properties"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Incorrect number of arguments: ", len(os.Args), ". Intended number - 2")
	}

	config := properties.GetConfig(os.Args[1])
	if err := apiServer.Start(config); err != nil {
		log.Fatal(err)
	}

}
