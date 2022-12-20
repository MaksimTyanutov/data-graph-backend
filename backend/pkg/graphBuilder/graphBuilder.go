package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/dbConnector"
	"data-graph-backend/pkg/properties"
	"errors"
)

var colors = []string{
	"#FF595E",
	"#FFCA3A",
	"#8AC926",
	"#1982C4",
	"#6A4C93",
}

// GET ALL PROJECTS
func GetProjects(dbConnector *dbConnector.PSQLConnector, minimized bool) ([]dataStructers.Project, error) {
	projects := make([]dataStructers.Project, 0)
	if !minimized {
		projectsDb, err := dbConnector.GetAllProjects()
		if err != nil {
			return nil, errors.New("graphBuilder:GetProjects(1). GetAllProjects don't work: " + err.Error())
		} else {
			for i := 0; i < len(projectsDb); i++ {
				project, err := projectsDb[i].Transform()
				if err != nil {
					return nil, errors.New("graphBuilder:GetProjects(2). Can't transform project. " + err.Error())
				}
				projects = append(projects, *project)
			}
		}
	} else {
		projectsDb, err := dbConnector.GetShortProjects()
		if err != nil {
			return nil, errors.New("graphBuilder:GetProjects(2). GetShortProjects don't work: " + err.Error())
		} else {
			for i := 0; i < len(projectsDb); i++ {
				project, err := projectsDb[i].Transform()
				if err != nil {
					return nil, errors.New("graphBuilder:GetProjects(2). Can't transform project. " + err.Error())
				}
				projects = append(projects, *project)
			}
		}
	}
	return projects, nil
}

// GET ALL COMPANIES
func GetCompanies(dbConnector *dbConnector.PSQLConnector) ([]dataStructers.Company, error) {
	companiesDb, err := dbConnector.GetAllCompanies()
	companies := make([]dataStructers.Company, 0)
	if err != nil {
		return nil, errors.New("graphBuilder:GetCompanies. GetAllCompanies don't work: " + err.Error())
	} else {
		for i := 0; i < len(companiesDb); i++ {
			company := companiesDb[i].Transform()
			companies = append(companies, company)
		}
	}
	return companies, nil
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
						Opacity: standardOpacity,
					})
				}
			} else {
				for j := 0; j < len(projects[i].PreviousNodeIds); j++ {
					isPresent := false
					for k := 0; k < len(projects); k++ {
						if projects[k].Id == projects[i].PreviousNodeIds[j] {
							isPresent = true
						}
					}
					if isPresent {
						links = append(links, Link{
							Source:  projects[i].PreviousNodeIds[j],
							Target:  projects[i].Id,
							Color:   colors[projects[i].CompanyId%len(colors)],
							Opacity: standardOpacity,
						})
					}
				}
				if projects[i-1].ProjectId == projects[i].ProjectId {
					links = append(links, Link{
						Source:  projects[i-1].Id,
						Target:  projects[i].Id,
						Color:   colors[projects[i].CompanyId%len(colors)],
						Opacity: standardOpacity,
					})
				}
			}
			if projects[i-1].CompanyId != projects[i].CompanyId {
				links = append(links, Link{
					Source:  projects[i].CompanyId + properties.CompanyIdShift,
					Target:  projects[i].Id,
					Color:   colors[projects[i].CompanyId%len(colors)],
					Opacity: standardOpacity,
				})
			}
		} else {
			links = append(links, Link{
				Source:  projects[i].CompanyId + properties.CompanyIdShift,
				Target:  projects[i].Id,
				Color:   colors[projects[i].CompanyId%len(colors)],
				Opacity: standardOpacity,
			})
		}
	}
	return links
}

func GetGraph(dbConnector *dbConnector.PSQLConnector, minimized bool) (*Graph, error) {
	err := dbConnector.SetIdShift()
	if err != nil {
		return nil, err
	}
	companies, err := GetCompanies(dbConnector)
	if err != nil {
		return nil, errors.New("GetGraph(2): " + err.Error())
	}
	projects, err := GetProjects(dbConnector, minimized)
	if err != nil {
		return nil, errors.New("GetGraph(3): " + err.Error())
	}
	nodes := TransformComp(companies)
	projectsTransformed := TransformProj(projects)
	nodes = append(nodes, projectsTransformed...)
	links := GetLinks(projects, minimized)
	return &Graph{
		Nodes: nodes,
		Links: links,
	}, nil
}
