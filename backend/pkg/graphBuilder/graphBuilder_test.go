package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"encoding/json"
	"testing"
)

func (l1 *Link) equal(l2 Link) bool {
	if l1.Source != l2.Source {
		return false
	}
	if l1.Target != l2.Target {
		return false
	}
	if l1.Color != l2.Color {
		return false
	}
	if l1.Opacity != l2.Opacity {
		return false
	}
	return true
}

func TestGetLinks(t *testing.T) {
	testTable := []struct {
		name          string
		projects      []dataStructers.Project
		expectedLinks []Link
		short         bool
	}{
		{
			name:          "Empty project",
			projects:      []dataStructers.Project{},
			expectedLinks: []Link{},
			short:         true,
		},
		{
			name: "Correct GetLinks short",
			projects: []dataStructers.Project{
				{
					Id:          2,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "Первая публичная версия поисковой системы Яндекс",
					Version:     "1.0",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date:            "1997-09-23T00:00:00Z",
					Url:             "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{},
					PressURL:        "",
					PosX:            600,
					PosY:            200,
				},
				{
					Id:          14,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "v5",
					Version:     "5.0",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date: "2011-08-19T00:00:00Z",
					Url:  "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{
						12,
					},
					PressURL: "",
					PosX:     600,
					PosY:     800,
				},
			},
			expectedLinks: []Link{
				{
					Source:  3,
					Target:  2,
					Color:   "#1982C4",
					Opacity: 1,
				},
				{
					Source:  2,
					Target:  14,
					Color:   "#1982C4",
					Opacity: 1,
				},
			},
			short: true,
		},
		{
			name: "Correct GetLinks full",
			projects: []dataStructers.Project{
				{
					Id:          2,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "Первая публичная версия поисковой системы Яндекс",
					Version:     "1.0",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date:            "1997-09-23T00:00:00Z",
					Url:             "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{},
					PressURL:        "",
					PosX:            600,
					PosY:            200,
				},
				{
					Id:          4,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "Яндекс Поисковик",
					Version:     "1.1",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date: "1998-10-16T00:00:00Z",
					Url:  "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{
						2,
					},
					PressURL: "",
					PosX:     600,
					PosY:     300,
				},
				{
					Id:          6,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "Следующая итерация",
					Version:     "1.2",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date: "1999-06-09T00:00:00Z",
					Url:  "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{
						4,
					},
					PressURL: "",
					PosX:     600,
					PosY:     400,
				},
				{
					Id:          8,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "Яндекс v2",
					Version:     "2.0",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date: "2001-06-13T00:00:00Z",
					Url:  "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{
						6,
					},
					PressURL: "",
					PosX:     600,
					PosY:     500,
				},
				{
					Id:          10,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "v3",
					Version:     "v3.0",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date: "2005-12-23T00:00:00Z",
					Url:  "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{
						8,
					},
					PressURL: "",
					PosX:     600,
					PosY:     600,
				},
				{
					Id:          12,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "v4",
					Version:     "4.0",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date: "2010-04-22T00:00:00Z",
					Url:  "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{
						10,
					},
					PressURL: "",
					PosX:     600,
					PosY:     700,
				},
				{
					Id:          14,
					ProjectId:   3,
					Name:        "Яndex-Web",
					Description: "v5",
					Version:     "5.0",
					CompanyId:   3,
					ProjectTypes: []string{
						"IT",
					},
					Date: "2011-08-19T00:00:00Z",
					Url:  "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					PreviousNodeIds: []int{
						12,
					},
					PressURL: "",
					PosX:     600,
					PosY:     800,
				},
			},
			expectedLinks: []Link{
				{
					Source:  3,
					Target:  2,
					Color:   "#1982C4",
					Opacity: 1,
				},
				{
					Source:  2,
					Target:  4,
					Color:   "#1982C4",
					Opacity: 1,
				},
				{
					Source:  4,
					Target:  6,
					Color:   "#1982C4",
					Opacity: 1,
				},
				{
					Source:  6,
					Target:  8,
					Color:   "#1982C4",
					Opacity: 1,
				},
				{
					Source:  8,
					Target:  10,
					Color:   "#1982C4",
					Opacity: 1,
				},
				{
					Source:  10,
					Target:  12,
					Color:   "#1982C4",
					Opacity: 1,
				},
				{
					Source:  12,
					Target:  14,
					Color:   "#1982C4",
					Opacity: 1,
				},
			},
			short: false,
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			links := GetLinks(test.projects, test.short)
			if links == nil && test.expectedLinks != nil {
				jsonStr, _ := json.Marshal(test.expectedLinks)
				t.Errorf("Actual str = nil\nExpected str = %s", jsonStr)
				return
			}
			if links != nil && test.expectedLinks == nil {
				jsonStr, _ := json.Marshal(links)
				t.Errorf("Actual str = %s\nExpected str = nil", jsonStr)
				return
			}
			if len(links) != len(test.expectedLinks) {
				t.Errorf("Actual len = %d\nExpected len = %d", len(links), len(test.expectedLinks))
				return
			} else {
				for i := 0; i < len(links); i++ {
					if links[i].equal(test.expectedLinks[i]) != true {
						jsonActual, _ := json.Marshal(links[i])
						jsonExpected, _ := json.Marshal(test.expectedLinks[i])
						t.Errorf("Actual link = %s\nExpected link = %s", jsonActual, jsonExpected)
						return
					}
				}
			}
		})
	}
}
