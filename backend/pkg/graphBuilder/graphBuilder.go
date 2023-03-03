package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"data-graph-backend/pkg/properties"
)

var colors = []string{
	"#FF595E",
	"#FFCA3A",
	"#8AC926",
	"#1982C4",
	"#6A4C93",
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
