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
	Logger      *logrus.Logger
	DbConnector *dbConnector.PSQLConnector
}

func ConfigureRouters(r *Router) {
	http.HandleFunc("/test", r.HandleTestAnswer)
	http.HandleFunc("/ping", r.HandlePingDB)
	http.HandleFunc("/Companies", r.handleCompanies)
	http.HandleFunc("/Projects", r.handleProjects)
	http.HandleFunc("/get:full", r.HandleGetGraphFull)
	http.HandleFunc("/get:short", r.HandleGetGraphShort)
	http.HandleFunc("/company", r.HandleCompany)
	http.HandleFunc("/product", r.HandleProduct)
	http.HandleFunc("/link/products", r.handleTimelineProduct)
	http.HandleFunc("/link/company", r.handleTimelineCompany)
	http.HandleFunc("/departments", r.handleGetAllDepartments)
	http.HandleFunc("/filterPresets", r.HandleGetFilterPresets)
	http.HandleFunc("/filterCompany", r.HandleFilterCompany)
	http.HandleFunc("/filterProduct", r.HandleFilterProduct)
}

func (rout *Router) HandleTestAnswer(rw http.ResponseWriter, r *http.Request) {
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, "test")
}

func (rout *Router) HandlePingDB(rw http.ResponseWriter, r *http.Request) {
	rout.setCorsHeaders(&rw)
	if err := rout.DbConnector.Ping(); err != nil {
		rout.respond(rw, r, http.StatusInternalServerError, "Database not responding")
		return
	}
	rout.respond(rw, r, http.StatusOK, "OK")
}

// GET ALL PROJECTS
func (rout *Router) handleProjects(rw http.ResponseWriter, r *http.Request) {
	projectsDb, err := rout.DbConnector.GetAllProjects()
	projects := make([]dataStructers.Project, 0)
	if err != nil {
		rout.Logger.Error("GetAllProjects don't work: " + err.Error())
	} else {
		for i := 0; i < len(projectsDb); i++ {
			project, err := projectsDb[i].Transform()
			if err != nil {
				rout.Logger.Error("GetAllProjects don't work: " + err.Error())
			}
			projects = append(projects, *project)
		}
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, projects)
}

// GET ALL COMPANIES
func (rout *Router) handleCompanies(rw http.ResponseWriter, r *http.Request) {
	companiesDb, err := rout.DbConnector.GetAllCompanies()
	companies := make([]dataStructers.Company, 0)
	if err != nil {
		rout.Logger.Error("GetAllCompanies don't work: ", err.Error())
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
func (rout *Router) HandleGetGraphFull(rw http.ResponseWriter, r *http.Request) {
	graph, err := rout.DbConnector.GetGraph(false)
	if err != nil {
		rout.Logger.Error("GetGraphFull don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, graph)
}

// GET GRAPH SHORT
func (rout *Router) HandleGetGraphShort(rw http.ResponseWriter, r *http.Request) {
	graph, err := rout.DbConnector.GetGraph(true)
	rout.Logger.Info("Sending short graph")
	if err != nil {
		rout.Logger.Error("GetGraphShort don't work: ", err.Error())
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, graph)
	rout.Logger.Info("Successful short graph")
}

// Get company information
func (rout *Router) HandleCompany(rw http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	rout.setCorsHeaders(&rw)
	companyID, err := strconv.Atoi(idStr)
	company, err := rout.DbConnector.GetCompanyInfo(companyID)
	if err != nil {
		rout.Logger.Error("GetCompanyInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, "Wrong argument: "+err.Error())
	}
	rout.respond(rw, r, http.StatusOK, company)
}

// Get product information
func (rout *Router) HandleProduct(rw http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	rout.setCorsHeaders(&rw)
	productID, err := strconv.Atoi(idStr)
	product, err := rout.DbConnector.GetProductInfo(productID)
	if err != nil {
		rout.Logger.Error("GetCompanyInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, "Wrong argument: "+err.Error())
	}
	rout.respond(rw, r, http.StatusOK, product)
}

// Get link information
func (rout *Router) handleTimelineCompany(rw http.ResponseWriter, r *http.Request) {
	targetIdStr := r.FormValue("target")
	targetID, err := strconv.Atoi(targetIdStr)
	if err != nil {
		rout.Logger.Error("Can't convert to number - "+targetIdStr+". Error: ", err.Error())
	}
	sourceIdStr := r.FormValue("source")
	sourceID, err := strconv.Atoi(sourceIdStr)
	if err != nil {
		rout.Logger.Error("Can't convert to number - "+sourceIdStr+". Error: ", err.Error())
	}

	source, err := rout.DbConnector.GetCompanyInfo(sourceID)
	if err != nil {
		rout.Logger.Error("GetCompanyInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	target, err := rout.DbConnector.GetProductInfo(targetID)
	if err != nil {
		rout.Logger.Error("GetProductInfo don't work: ", err.Error())
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
		rout.Logger.Error("Can't convert to number - "+targetIdStr+". Error: ", err.Error())
	}
	sourceIdStr := r.FormValue("source")
	sourceID, err := strconv.Atoi(sourceIdStr)
	if err != nil {
		rout.Logger.Error("Can't convert to number - "+sourceIdStr+". Error: ", err.Error())
	}

	source, err := rout.DbConnector.GetProductInfo(sourceID)
	if err != nil {
		rout.Logger.Error("GetProductInfo don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	target, err := rout.DbConnector.GetProductInfo(targetID)
	if err != nil {
		rout.Logger.Error("GetProductInfo don't work: ", err.Error())
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
	departments, err := rout.DbConnector.GetAllDepartments()
	if err != nil {
		rout.Logger.Error("GetAllDepartments don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, departments)
}

func (rout *Router) HandleGetFilterPresets(rw http.ResponseWriter, r *http.Request) {
	companyFilters, err := rout.DbConnector.GetCompanyFilters()
	if err != nil {
		rout.Logger.Error("GetFilterPresets(1) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	productFilters, err := rout.DbConnector.GetProductFilters()
	if err != nil {
		rout.Logger.Error("GetFilterPresets(2) don't work: ", err.Error())
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

func (rout *Router) HandleFilterCompany(rw http.ResponseWriter, r *http.Request) {
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
		rout.Logger.Error("FilterCompany(1) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(data, &companyFilters)
	if err != nil {
		rout.Logger.Error("FilterCompany(2). Filters unmarshall don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = utils.ValidateFilterCompany(companyFilters)
	if err != nil {
		rout.Logger.Error(err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	idArray, err := rout.DbConnector.GetFiltersIDCompany(companyFilters)
	if err != nil {
		rout.Logger.Error("GetFiltersID don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	rout.setCorsHeaders(&rw)
	rout.respond(rw, r, http.StatusOK, idArray)
}

func (rout *Router) HandleFilterProduct(rw http.ResponseWriter, r *http.Request) {
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
		rout.Logger.Error("HandleFilterProduct(1) don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(data, &productFilters)
	if err != nil {
		rout.Logger.Error("HandleFilterProduct(2). Filters unmarshall don't work: ", err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	err = utils.ValidateFilterProduct(productFilters)
	if err != nil {
		rout.Logger.Error(err.Error())
		rout.respond(rw, r, http.StatusBadRequest, err)
		return
	}
	idArray, err := rout.DbConnector.GetFiltersIDProduct(productFilters)
	if err != nil {
		rout.Logger.Error("HandleFilterProduct(3) don't work: ", err.Error())
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
			rout.Logger.Error("Error while responding: " + err.Error() + ".\nRequest: " + r.URL.String())
		}
	}
}
