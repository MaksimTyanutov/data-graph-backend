package apiServer

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/graphBuilder"
	"encoding/json"
	"log"
	"net/http"
)

type Router struct {
	dbConnector *dbConnector.PSQLConnector
}

func configureRouters(r *Router) {
	http.HandleFunc("/test", r.handleTestAnswer)
	http.HandleFunc("/Companies", r.handleCompanies)
	http.HandleFunc("/Projects", r.handleProjects)
	http.HandleFunc("/get:full", r.handleGetGraphFull)
	http.HandleFunc("/get:short", r.handleGetGraphShort)
}

func (rout *Router) handleTestAnswer(rw http.ResponseWriter, r *http.Request) {
	respond(rw, r, http.StatusOK, "test")
}

// GET ALL PROJECTS
func (rout *Router) handleProjects(rw http.ResponseWriter, r *http.Request) {
	projectsDb, err := rout.dbConnector.GetAllProjects()
	projects := make([]dataStructers.Project, 0)
	if err != nil {
		log.Print("GetAllProjects don't work: ", err.Error())
	} else {
		for i := 0; i < len(projectsDb); i++ {
			project := projectsDb[i].Transform()
			projects = append(projects, project)
		}
	}
	respond(rw, r, http.StatusOK, projects)
}

// GET ALL COMPANIES
func (rout *Router) handleCompanies(rw http.ResponseWriter, r *http.Request) {
	companiesDb, err := rout.dbConnector.GetAllCompanies()
	companies := make([]dataStructers.Company, 0)
	if err != nil {
		log.Print("GetAllCompanies don't work: ", err.Error())
	} else {
		for i := 0; i < len(companiesDb); i++ {
			company := companiesDb[i].Transform()
			companies = append(companies, company)
		}
	}
	respond(rw, r, http.StatusOK, companies)
}

// GET GRAPH FULL
func (rout *Router) handleGetGraphFull(rw http.ResponseWriter, r *http.Request) {
	graph := graphBuilder.GetGraph(rout.dbConnector, false)
	respond(rw, r, http.StatusOK, graph)
}

// GET GRAPH SHORT
func (rout *Router) handleGetGraphShort(rw http.ResponseWriter, r *http.Request) {
	graph := graphBuilder.GetGraph(rout.dbConnector, true)
	respond(rw, r, http.StatusOK, graph)
}

//func parseError(w http.ResponseWriter, r *http.Request, code int, err error) {
//	respond(w, r, code, map[string]string{"error": err.Error()})
//}

func respond(w http.ResponseWriter, r *http.Request, code int, date interface{}) {
	w.WriteHeader(code)
	if date != nil {
		err := json.NewEncoder(w).Encode(date)
		if err != nil {
			log.Print("Error while responding: " + err.Error() + ".\nRequest: " + r.URL.String())
		}
	}
}
