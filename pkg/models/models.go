package models

import (
	"strconv"
	"time"
)

type Project struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	FileName string    `json:"-"`
	Created  time.Time `json:"created"`
}

func (p *Project) FilePath() string {
	if p.FileName == "" {
		return "/static/img/project_placeholder.png"
	}
	return "/store/" + strconv.Itoa(p.ID) + "/" + p.FileName
}

type Projects []*Project

type Photo struct {
	ID            int
	ProjectID     int
	RunningNumber int
	FileName      string
	Location      string
	Created       time.Time
}

type Photos []*Photo

func (f *Photo) FilePath() string {
	return "/store/" + strconv.Itoa(f.ProjectID) + "/" + f.FileName
}

type FormField struct {
	ID    int
	Name  string
	Title string
}

type FormFields []*FormField
