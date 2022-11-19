package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"log"
)

// GET ALL PROJECTS
func GetProjects(dbConnector *dbConnector.PSQLConnector) []dataStructers.Project {
	projectsDb, err := dbConnector.GetAllProjects()
	projects := make([]dataStructers.Project, 0)
	if err != nil {
		log.Print("GetAllProjects don't work: ", err.Error())
	} else {
		for i := 0; i < len(projectsDb); i++ {
			project := projectsDb[i].Transform()
			projects = append(projects, project)
		}
	}
	return projects
}

// GET ALL COMPANIES
func GetCompanies(dbConnector *dbConnector.PSQLConnector) []dataStructers.Company {
	companiesDb, err := dbConnector.GetAllCompanies()
	companies := make([]dataStructers.Company, 0)
	if err != nil {
		log.Print("GetAllCompanies don't work: ", err.Error())
	} else {
		for i := 0; i < len(companiesDb); i++ {
			company := companiesDb[i].Transform()
			companies = append(companies, company)
		}
	}
	return companies
}

func GetLinks(projects []dataStructers.Project) []Link {
	links := make([]Link, 0)
	for i := 0; i < len(projects); i++ {
		if len(projects[i].PreviousNodeIds) != 0 {
			for j := 0; j < len(projects[i].PreviousNodeIds); j++ {
				links = append(links, Link{
					Source: projects[i].PreviousNodeIds[j],
					Target: projects[i].Id,
				})
			}
		} else {
			links = append(links, Link{
				Source: projects[i].CompanyId + companyIdShift,
				Target: projects[i].Id,
			})
		}
	}
	return links
}

func GetGraph(dbConnector *dbConnector.PSQLConnector) Graph {
	companies := GetCompanies(dbConnector)
	projects := GetProjects(dbConnector)
	nodes := TransformComp(companies)
	nodes = append(nodes, TransformProj(projects)...)
	links := GetLinks(projects)
	return Graph{
		Nodes: nodes,
		Links: links,
	}
}
