package apiServer

import (
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/logging"
	"data-graph-backend/pkg/properties"
	"log"
	"net/http"
)

func Start(config *properties.Config) error {

	logger := logging.GetLogger()
	logger.Info("Starting apiServer")

	logger.Info("Trying to connect to DB")
	dbConnection, err := dbConnector.NewConnection(config)
	if err != nil {
		logger.Fatal("Can't connect to db - ", err.Error())
	}

	logger.Info("Connection successful")
	router := &Router{
		logger:      logger,
		dbConnector: dbConnection,
	}

	configureRouters(router)

	logger.Info("Listening on " + config.ProgramSettings.Host + config.ProgramSettings.Port)
	log.Println("Listening on " + config.ProgramSettings.Host + config.ProgramSettings.Port)
	return http.ListenAndServe(config.ProgramSettings.Port, nil)
}
