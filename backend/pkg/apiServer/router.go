package apiServer

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/utils"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type Router struct {
	logger      *logrus.Logger
	dbConnector *dbConnector.PSQLConnector
}

func configureRouters(r *Router) {
	http.HandleFunc("/test", r.handleTestAnswer)
	http.HandleFunc("/ping", r.handlePing)
	http.HandleFunc("/Companies", r.handleCompanies)
	http.HandleFunc("/Projects", r.handleProjects)
	http.HandleFunc("/get:full", r.handleGetGraphFull)
	http.HandleFunc("/get:short", r.handleGetGraphShort)
	http.HandleFunc("/company", r.handleCompany)
	http.HandleFunc("/product", r.handleProduct)
	http.HandleFunc("/link/products", r.handleTimelineProduct)
	http.HandleFunc("/link/company", r.handleTimelineCompany)
	http.HandleFunc("/departments", r.handleGetAllDepartments)
	http.HandleFunc("/filterPresets", r.handleGetFilterPresets)
	http.HandleFunc("/filterCompany", r.handleFilterCompany)
	http.HandleFunc("/filterProduct", r.handleFilterProduct)
}

func (rout *Router) handleTestAnswer(rw http.ResponseWriter, r *http.Request) {
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, "test")
}

func (rout *Router) handlePing(rw http.ResponseWriter, r *http.Request) {
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, "OK")
}

// GET ALL PROJECTS
func (rout *Router) handleProjects(rw http.ResponseWriter, r *http.Request) {
	projectsDb, err := rout.dbConnector.GetAllProjects()
	projects := make([]dataStructers.Project, 0)
	if err != nil {
		rout.logger.Error("GetAllProjects don't work: " + err.Error())
	} else {
		for i := 0; i < len(projectsDb); i++ {
			project, err := projectsDb[i].Transform()
			if err != nil {
				rout.logger.Error("GetAllProjects don't work: " + err.Error())
			}
			projects = append(projects, *project)
		}
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, projects)
}

// GET ALL COMPANIES
func (rout *Router) handleCompanies(rw http.ResponseWriter, r *http.Request) {
	companiesDb, err := rout.dbConnector.GetAllCompanies()
	companies := make([]dataStructers.Company, 0)
	if err != nil {
		rout.logger.Error("GetAllCompanies don't work: ", err.Error())
	} else {
		for i := 0; i < len(companiesDb); i++ {
			company := companiesDb[i].Transform()
			companies = append(companies, company)
		}
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, companies)
}

// GET GRAPH FULL
func (rout *Router) handleGetGraphFull(rw http.ResponseWriter, r *http.Request) {
	graph, err := rout.dbConnector.GetGraph(false)
	if err != nil {
		rout.logger.Error("GetGraphFull don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, graph)
}

// GET GRAPH SHORT
func (rout *Router) handleGetGraphShort(rw http.ResponseWriter, r *http.Request) {
	graph, err := rout.dbConnector.GetGraph(true)
	rout.logger.Info("Sending short graph")
	if err != nil {
		rout.logger.Error("GetGraphShort don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, graph)
	rout.logger.Info("Successful short graph")
}

// Get company information
func (rout *Router) handleCompany(rw http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	companyID, err := strconv.Atoi(idStr)
	company, err := rout.dbConnector.GetCompanyInfo(companyID)
	if err != nil {
		rout.logger.Error("GetCompanyInfo don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, company)
}

// Get product information
func (rout *Router) handleProduct(rw http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	productID, err := strconv.Atoi(idStr)
	product, err := rout.dbConnector.GetProductInfo(productID)
	if err != nil {
		rout.logger.Error("GetCompanyInfo don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, product)
}

// Get link information
func (rout *Router) handleTimelineCompany(rw http.ResponseWriter, r *http.Request) {
	targetIdStr := r.FormValue("target")
	targetID, err := strconv.Atoi(targetIdStr)
	if err != nil {
		rout.logger.Error("Can't convert to number - "+targetIdStr+". Error: ", err.Error())
	}
	sourceIdStr := r.FormValue("source")
	sourceID, err := strconv.Atoi(sourceIdStr)
	if err != nil {
		rout.logger.Error("Can't convert to number - "+sourceIdStr+". Error: ", err.Error())
	}

	source, err := rout.dbConnector.GetCompanyInfo(sourceID)
	if err != nil {
		rout.logger.Error("GetCompanyInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	target, err := rout.dbConnector.GetProductInfo(targetID)
	if err != nil {
		rout.logger.Error("GetProductInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	timeline := dataStructers.TimelineCompany{
		Company: struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Year string `json:"year"`
		}{
			Id:   source.Id,
			Name: source.Name,
			Year: source.FoundationYear,
		},
		Product: struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Year string `json:"year"`
		}{
			Id:   target.Id,
			Name: target.Name,
			Year: target.Year,
		},
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, timeline)
}

// Get link information
func (rout *Router) handleTimelineProduct(rw http.ResponseWriter, r *http.Request) {
	targetIdStr := r.FormValue("target")
	targetID, err := strconv.Atoi(targetIdStr)
	if err != nil {
		rout.logger.Error("Can't convert to number - "+targetIdStr+". Error: ", err.Error())
	}
	sourceIdStr := r.FormValue("source")
	sourceID, err := strconv.Atoi(sourceIdStr)
	if err != nil {
		rout.logger.Error("Can't convert to number - "+sourceIdStr+". Error: ", err.Error())
	}

	source, err := rout.dbConnector.GetProductInfo(sourceID)
	if err != nil {
		rout.logger.Error("GetProductInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	target, err := rout.dbConnector.GetProductInfo(targetID)
	if err != nil {
		rout.logger.Error("GetProductInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	timeline := dataStructers.TimelineProducts{
		Product1: struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Year string `json:"year"`
		}{
			Id:   source.Id,
			Name: source.Name,
			Year: source.Year,
		},
		Product2: struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Year string `json:"year"`
		}{
			Id:   target.Id,
			Name: target.Name,
			Year: target.Year,
		},
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, timeline)
}

func (rout *Router) handleGetAllDepartments(rw http.ResponseWriter, r *http.Request) {
	departments, err := rout.dbConnector.GetAllDepartments()
	if err != nil {
		rout.logger.Error("GetAllDepartments don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, departments)
}

func (rout *Router) handleGetFilterPresets(rw http.ResponseWriter, r *http.Request) {
	companyFilters, err := rout.dbConnector.GetCompanyFilters()
	if err != nil {
		rout.logger.Error("GetFilterPresets(1) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	productFilters, err := rout.dbConnector.GetProductFilters()
	if err != nil {
		rout.logger.Error("GetFilterPresets(2) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	filterPresets := dataStructers.FilterPresets{
		CompanyFilters: *companyFilters,
		ProductFilters: *productFilters,
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, filterPresets)
}

func (rout *Router) handleFilterCompany(rw http.ResponseWriter, r *http.Request) {
	rout.setCorsHeaders(&rw)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		rout.respond(rw, r, http.StatusMethodNotAllowed, nil)
		return
	}

	var companyFilters dataStructers.CompanyFilters
	data, err := io.ReadAll(r.Body)
	if err != nil {
		rout.logger.Error("FilterCompany(1) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(data, &companyFilters)
	if err != nil {
		rout.logger.Error("FilterCompany(2). Filters unmarshall don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = utils.ValidateFilterCompany(companyFilters)
	if err != nil {
		rout.logger.Error(err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	idArray, err := rout.dbConnector.GetFiltersIDCompany(companyFilters)
	if err != nil {
		rout.logger.Error("GetFiltersID don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, idArray)
}

func (rout *Router) handleFilterProduct(rw http.ResponseWriter, r *http.Request) {
	rout.setCorsHeaders(&rw)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != http.MethodPost {
		rout.respond(rw, r, http.StatusMethodNotAllowed, nil)
		return
	}

	var productFilters dataStructers.ProductFilters
	data, err := io.ReadAll(r.Body)
	if err != nil {
		rout.logger.Error("handleFilterProduct(1) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(data, &productFilters)
	if err != nil {
		rout.logger.Error("handleFilterProduct(2). Filters unmarshall don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = utils.ValidateFilterProduct(productFilters)
	if err != nil {
		rout.logger.Error(err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	idArray, err := rout.dbConnector.GetFiltersIDProduct(productFilters)
	if err != nil {
		rout.logger.Error("handleFilterProduct(3) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, idArray)
}

//func parseError(w http.ResponseWriter, r *http.Request, code int, err error) {
//	respond(w, r, code, map[string]string{"error": err.Error()})
//}

func (rout *Router) setCorsHeaders(rw *http.ResponseWriter) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "*")
	(*rw).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (rout *Router) respond(w http.ResponseWriter, r *http.Request, code int, date interface{}) {
	w.WriteHeader(code)

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	// return "OKOK"
	if date != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := json.NewEncoder(w).Encode(date)
		if err != nil {
			rout.logger.Error("Error while responding: " + err.Error() + ".\nRequest: " + r.URL.String())
		}
	}
}
