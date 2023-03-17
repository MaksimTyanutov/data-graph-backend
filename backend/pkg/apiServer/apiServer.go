package apiServer

import (
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/properties"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func Start(config *properties.Config, logger *logrus.Logger) error {

	logger.Info("Starting apiServer")

	logger.Info("Trying to connect to DB")
	dbConnection, err := dbConnector.NewConnection(config, logger)
	if err != nil {
		logger.Fatal("Can't connect to db - ", err.Error())
	}
	err = dbConnection.SetIdShift()
	if err != nil {
		logger.Fatal("Can't get info from db - ", err.Error())
	}

	logger.Info("Connection successful")
	router := &Router{
		Logger:      logger,
		DbConnector: dbConnection,
	}

	ConfigureRouters(router)

	logger.Info("Listening on " + config.ProgramSettings.Host + config.ProgramSettings.Port)
	log.Println("Listening on " + config.ProgramSettings.Host + config.ProgramSettings.Port)
	return http.ListenAndServe(config.ProgramSettings.Port, nil)
}
