package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"log"
)

var colors = []string{
	"808080",
	"FFFFFF",
	"800000",
	"FF0000",
	"800080",
	"FF00FF",
	"008000",
	"00FF00",
	"808000",
	"FFFF00",
	"000080",
	"0000FF",
	"008080",
	"00FFFF",
}

// GET ALL PROJECTS
func GetProjects(dbConnector *dbConnector.PSQLConnector, minimized bool) []dataStructers.Project {
	projects := make([]dataStructers.Project, 0)
	if !minimized {
		projectsDb, err := dbConnector.GetAllProjects()
		if err != nil {
			log.Print("GetAllProjects don't work: ", err.Error())
		} else {
			for i := 0; i < len(projectsDb); i++ {
				project := projectsDb[i].Transform()
				projects = append(projects, project)
			}
		}
	} else {
		projectsDb, err := dbConnector.GetShortProjects()
		if err != nil {
			log.Print("GetShortProjects don't work: ", err.Error())
		} else {
			for i := 0; i < len(projectsDb); i++ {
				project := projectsDb[i].Transform()
				projects = append(projects, project)
			}
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
					Color:  colors[projects[i].CompanyId%len(colors)],
				})
			}
		} else {
			links = append(links, Link{
				Source: projects[i].CompanyId + companyIdShift,
				Target: projects[i].Id,
				Color:  colors[projects[i].CompanyId%len(colors)],
			})
		}
	}
	return links
}

func GetGraph(dbConnector *dbConnector.PSQLConnector, minimized bool) Graph {
	companies := GetCompanies(dbConnector)
	projects := GetProjects(dbConnector, minimized)
	nodes := TransformComp(companies)
	nodes = append(nodes, TransformProj(projects)...)
	links := GetLinks(projects)
	return Graph{
		Nodes: nodes,
		Links: links,
	}
}
