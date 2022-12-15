package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/properties"
)

var sizeCompany = 900
var sizeProduct = 600
var standartOpacity float32 = 1.0

type Node struct {
	Id       int     `json:"id"`
	Svg      string  `json:"svg"`
	NodeType string  `json:"nodeType"`
	Size     int     `json:"size"`
	Opacity  float32 `json:"opacity"`
}

type Link struct {
	Source  int     `json:"source"`
	Target  int     `json:"target"`
	Color   string  `json:"color"`
	Opacity float32 `json:"opacity"`
}

func TransformComp(companies []dataStructers.Company) []Node {
	nodes := make([]Node, 0)
	for i := 0; i < len(companies); i++ {
		node := Node{
			Id:       companies[i].Id + properties.CompanyIdShift,
			Svg:      companies[i].IconPath,
			NodeType: "Компания",
			Size:     sizeCompany,
			Opacity:  standartOpacity,
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
			Size:     sizeProduct,
			Opacity:  standartOpacity,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}
