package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/properties"
	"log"
)

var colors = []string{
	"#808080",
	"#000000",
	"#FF0000",
	"#800080",
	"#008000",
	"#00FF00",
	"#808000",
	"#800000",
	"#FFFF00",
	"#000080",
	"#0000FF",
	"#008080",
	"#00FFFF",
}

// GET ALL PROJECTS
func GetProjects(dbConnector *dbConnector.PSQLConnector, minimized bool) []dataStructers.Project {
	projects := make([]dataStructers.Project, 0)
	if !minimized {
		projectsDb, err := dbConnector.GetAllProjects()
		if err != nil {
			log.Print("graphBuilder:GetProjects(1). GetAllProjects don't work: ", err.Error())
		} else {
			for i := 0; i < len(projectsDb); i++ {
				project := projectsDb[i].Transform()
				projects = append(projects, project)
			}
		}
	} else {
		projectsDb, err := dbConnector.GetShortProjects()
		if err != nil {
			log.Print("graphBuilder:GetProjects(2). GetShortProjects don't work: ", err.Error())
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
		log.Print("graphBuilder:GetCompanies. GetAllCompanies don't work: ", err.Error())
	} else {
		for i := 0; i < len(companiesDb); i++ {
			company := companiesDb[i].Transform()
			companies = append(companies, company)
		}
	}
	return companies
}

func GetLinks(projects []dataStructers.Project, short bool) []Link {
	links := make([]Link, 0)
	for i := 0; i < len(projects); i++ {
		if len(projects[i].PreviousNodeIds) != 0 {
			if short == false {
				for j := 0; j < len(projects[i].PreviousNodeIds); j++ {
					links = append(links, Link{
						Source:  projects[i].PreviousNodeIds[j],
						Target:  projects[i].Id,
						Color:   colors[projects[i].CompanyId%len(colors)],
						Opacity: standartOpacity,
					})
				}
			} else {
				for j := 0; j < len(projects[i].PreviousNodeIds)-1; j++ {
					links = append(links, Link{
						Source:  projects[i].PreviousNodeIds[j],
						Target:  projects[i].Id,
						Color:   colors[projects[i].CompanyId%len(colors)],
						Opacity: standartOpacity,
					})
				}
				links = append(links, Link{
					Source:  projects[i-1].Id,
					Target:  projects[i].Id,
					Color:   colors[projects[i].CompanyId%len(colors)],
					Opacity: standartOpacity,
				})
			}
		} else {
			links = append(links, Link{
				Source:  projects[i].CompanyId + properties.CompanyIdShift,
				Target:  projects[i].Id,
				Color:   colors[projects[i].CompanyId%len(colors)],
				Opacity: standartOpacity,
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
	links := GetLinks(projects, minimized)
	return Graph{
		Nodes: nodes,
		Links: links,
	}
}
