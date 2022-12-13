package apiServer

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/graphBuilder"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Router struct {
	dbConnector *dbConnector.PSQLConnector
}

func configureRouters(r *Router) {
	http.HandleFunc("/test", r.handleTestAnswer)
	//http.HandleFunc("/Companies", r.handleCompanies)
	//http.HandleFunc("/Projects", r.handleProjects)
	http.HandleFunc("/get:full", r.handleGetGraphFull)
	http.HandleFunc("/get:short", r.handleGetGraphShort)
	http.HandleFunc("/company", r.handleCompany)
	http.HandleFunc("/product", r.handleProduct)
}

func (rout *Router) handleTestAnswer(rw http.ResponseWriter, r *http.Request) {
	rout.setCorsHeaders(&rw)
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
	rout.setCorsHeaders(&rw)
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
	rout.setCorsHeaders(&rw)
	respond(rw, r, http.StatusOK, companies)
}

// GET GRAPH FULL
func (rout *Router) handleGetGraphFull(rw http.ResponseWriter, r *http.Request) {
	graph := graphBuilder.GetGraph(rout.dbConnector, false)
	rout.setCorsHeaders(&rw)
	respond(rw, r, http.StatusOK, graph)
}

// GET GRAPH SHORT
func (rout *Router) handleGetGraphShort(rw http.ResponseWriter, r *http.Request) {
	graph := graphBuilder.GetGraph(rout.dbConnector, true)
	rout.setCorsHeaders(&rw)
	respond(rw, r, http.StatusOK, graph)
}

func (rout *Router) setCorsHeaders(rw *http.ResponseWriter) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "*")
	(*rw).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Get company information
func (rout *Router) handleCompany(rw http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	companyID, err := strconv.Atoi(idStr)
	company, err := rout.dbConnector.GetCompanyInfo(companyID)
	if err != nil {
		log.Print("GetCompanyInfo don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	respond(rw, r, http.StatusOK, company)
}

// Get product information
func (rout *Router) handleProduct(rw http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	productID, err := strconv.Atoi(idStr)
	product, err := rout.dbConnector.GetProductInfo(productID)
	if err != nil {
		log.Print("GetCompanyInfo don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	respond(rw, r, http.StatusOK, product)
}

//func parseError(w http.ResponseWriter, r *http.Request, code int, err error) {
//	respond(w, r, code, map[string]string{"error": err.Error()})
//}

func respond(w http.ResponseWriter, r *http.Request, code int, date interface{}) {
	w.WriteHeader(code)

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// return "OKOK"
	if date != nil {
		err := json.NewEncoder(w).Encode(date)
		if err != nil {
			log.Print("Error while responding: " + err.Error() + ".\nRequest: " + r.URL.String())
		}
	}
}
