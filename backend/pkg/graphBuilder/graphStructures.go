package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/properties"
)

var sizeCompany = 900
var sizeProduct = 600
var standardOpacity float32 = 1.0

type Node struct {
	Name     string  `json:"name"`
	Id       int     `json:"id"`
	Svg      string  `json:"svg"`
	NodeType string  `json:"nodeType"`
	Size     int     `json:"size"`
	Opacity  float32 `json:"opacity"`
	Color    string  `json:"color"`
	X        int     `json:"x"`
	Y        int     `json:"y"`
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
			Name:     companies[i].Name,
			Id:       companies[i].Id + properties.CompanyIdShift,
			Svg:      companies[i].IconPath,
			NodeType: "Компания",
			Size:     sizeCompany,
			Opacity:  standardOpacity,
			Color:    colors[companies[i].Id%len(colors)],
			X:        companies[i].PosX,
			Y:        companies[i].PosY,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func TransformProj(projects []dataStructers.Project) []Node {
	nodes := make([]Node, 0)
	for i := 0; i < len(projects); i++ {
		node := Node{
			Name:     projects[i].Name,
			Id:       projects[i].Id,
			Svg:      projects[i].Url,
			NodeType: "Продукт",
			Size:     sizeProduct,
			Opacity:  standardOpacity,
			Color:    colors[projects[i].CompanyId%len(colors)],
			X:        projects[i].PosX,
			Y:        projects[i].PosY,
		}
		nodes = append(nodes, node)
	}
	return nodes
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}
