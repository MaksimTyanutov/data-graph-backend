package graphBuilder

import (
	"data-graph-backend/pkg/dataStructers"
	"testing"
)

func TestTransformComp(t *testing.T) {
	testTable := []struct {
		name          string
		Companies     []dataStructers.Company
		expectedNodes []Node
	}{
		{
			name: "Correct transform of one company",
			Companies: []dataStructers.Company{
				{Id: 2,
					Name:           "VK",
					Description:    "VK — крупнейшая российская технологическая компания, которая объединяет социальные сети, коммуникационные сервисы, игры, образование и многое другое. Проектами VK пользуются больше 90% аудитории рунета. Сегодня VK это не только социальная сеть, но и более 200 продуктов и сервисов для людей и бизнеса.",
					EmployeeNum:    8842,
					FoundationYear: "1998-10-01T00:00:00Z",
					CompanyTypeName: []string{
						"Связь",
						"Финансы",
						"IT",
						"Торговля",
						"Развлечения",
					},
					OwnerName: "Владимир Сергеевич Кириенко",
					Address:   "Санкт-Петербург",
					IconPath:  "http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png",
					PosX:      100,
					PosY:      0,
				},
			},
			expectedNodes: []Node{
				{
					Name:     "VK",
					Id:       79,
					Svg:      "http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png",
					NodeType: "Компания",
					Size:     900,
					Opacity:  1,
					Color:    "#8AC926",
					X:        100,
					Y:        0,
				},
			},
		},
		{
			name: "Correct transform of several company",
			Companies: []dataStructers.Company{
				{Id: 2,
					Name:           "VK",
					Description:    "VK — крупнейшая российская технологическая компания, которая объединяет социальные сети, коммуникационные сервисы, игры, образование и многое другое. Проектами VK пользуются больше 90% аудитории рунета. Сегодня VK это не только социальная сеть, но и более 200 продуктов и сервисов для людей и бизнеса.",
					EmployeeNum:    8842,
					FoundationYear: "1998-10-01T00:00:00Z",
					CompanyTypeName: []string{
						"Связь",
						"Финансы",
						"IT",
						"Торговля",
						"Развлечения",
					},
					OwnerName: "Владимир Сергеевич Кириенко",
					Address:   "Санкт-Петербург",
					IconPath:  "http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png",
					PosX:      100,
					PosY:      0,
				},
				{Id: 2,
					Name:           "VK",
					Description:    "VK — крупнейшая российская технологическая компания, которая объединяет социальные сети, коммуникационные сервисы, игры, образование и многое другое. Проектами VK пользуются больше 90% аудитории рунета. Сегодня VK это не только социальная сеть, но и более 200 продуктов и сервисов для людей и бизнеса.",
					EmployeeNum:    8842,
					FoundationYear: "1998-10-01T00:00:00Z",
					CompanyTypeName: []string{
						"Связь",
						"Финансы",
						"IT",
						"Торговля",
						"Развлечения",
					},
					OwnerName: "Владимир Сергеевич Кириенко",
					Address:   "Санкт-Петербург",
					IconPath:  "http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png",
					PosX:      100,
					PosY:      0,
				},
			},
			expectedNodes: []Node{
				{
					Name:     "VK",
					Id:       79,
					Svg:      "http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png",
					NodeType: "Компания",
					Size:     900,
					Opacity:  1,
					Color:    "#8AC926",
					X:        100,
					Y:        0,
				},
				{
					Name:     "VK",
					Id:       79,
					Svg:      "http://141.95.127.215:8082/media-server/data/f65df5df0cb6241377deafb69d35d620.png",
					NodeType: "Компания",
					Size:     900,
					Opacity:  1,
					Color:    "#8AC926",
					X:        100,
					Y:        0,
				},
			},
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			comp := TransformComp(test.Companies)
			if len(comp) != len(test.expectedNodes) {
				t.Errorf("Actual len = %d\nExpected len = %d", len(comp), len(test.expectedNodes))
				return
			}
			if len(comp) != 0 {
				for i := 0; i < len(comp); i++ {
					if comp[i].Name != test.expectedNodes[i].Name {
						t.Errorf("Actual Name = %d\nExpected Name = %d", len(comp[i].Name), len(test.expectedNodes[i].Name))
						return
					}
				}
			}
		})
	}
}

func TestTransformProj(t *testing.T) {
	testTable := []struct {
		name          string
		Projects      []dataStructers.Project
		expectedNodes []Node
	}{
		{
			name: "Correct transform of one project",
			Projects: []dataStructers.Project{
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
			},
			expectedNodes: []Node{
				{
					Name:     "Яndex-Web",
					Id:       2,
					Svg:      "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					NodeType: "Продукт",
					Size:     600,
					Opacity:  1,
					Color:    "#1982C4",
					X:        600,
					Y:        200,
				},
			},
		},
		{
			name: "Correct transform of several company",
			Projects: []dataStructers.Project{
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
			},
			expectedNodes: []Node{
				{
					Name:     "Яndex-Web",
					Id:       2,
					Svg:      "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					NodeType: "Продукт",
					Size:     600,
					Opacity:  1,
					Color:    "#1982C4",
					X:        600,
					Y:        200,
				}, {
					Name:     "Яndex-Web",
					Id:       2,
					Svg:      "http://141.95.127.215:8082/media-server/data/5a72230e4285ab8545ed460d662ac070.png",
					NodeType: "Продукт",
					Size:     600,
					Opacity:  1,
					Color:    "#1982C4",
					X:        600,
					Y:        200,
				},
			},
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			proj := TransformProj(test.Projects)
			if len(proj) != len(test.expectedNodes) {
				t.Errorf("Actual len = %d\nExpected len = %d", len(proj), len(test.expectedNodes))
				return
			}
			if len(proj) != 0 {
				for i := 0; i < len(proj); i++ {
					if proj[i].Name != test.expectedNodes[i].Name {
						t.Errorf("Actual Name = %d\nExpected Name = %d", len(proj[i].Name), len(test.expectedNodes[i].Name))
						return
					}
				}
			}
		})
	}
}
