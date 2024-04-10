package models

import "github.com/warmans/rsk-search/gen/api"

type CommunityProjects []CommunityProject

func (c CommunityProjects) Proto(totalCount int64) *api.CommunityProjectList {
	list := &api.CommunityProjectList{
		Projects:    nil,
		ResultCount: int32(totalCount),
	}
	for _, v := range c {
		list.Projects = append(list.Projects, v.Proto())
	}
	return list
}

type CommunityProject struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

func (c CommunityProject) Proto() *api.CommunityProject {
	return &api.CommunityProject{
		Id:        c.ID,
		Name:      c.Name,
		Summary:   c.Summary,
		Content:   c.Content,
		Url:       c.URL,
		CreatedAt: c.CreatedAt,
	}
}
