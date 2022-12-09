package graphBuilder

import "data-graph-backend/pkg/dataStructers"

var companyIdShift = 100000

type Node struct {
	Id       int    `json:"id"`
	Svg      string `json:"svg"`
	NodeType string `json:"nodeType"`
}

type Link struct {
	Source int    `'json:"source"`
	Target int    `'json:"target"`
	Color  string `json:"color"`
}

func TransformComp(companies []dataStructers.Company) []Node {
	nodes := make([]Node, 0)
	for i := 0; i < len(companies); i++ {
		node := Node{
			Id:       companies[i].Id + companyIdShift,
			Svg:      companies[i].IconPath,
			NodeType: "Компания",
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func TransformProj(projects []dataStructers.Project) []Node {
	nodes := make([]Node, 0)
	for i := 0; i < len(projects); i++ {
		node := Node{
			Id:       projects[i].Id,
			Svg:      "",
			NodeType: "Продукт",
		}
		nodes = append(nodes, node)
	}
	return nodes
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}
