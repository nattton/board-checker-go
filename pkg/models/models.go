package models

import (
	"strconv"
	"time"
)

type Worksheet struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type Worksheets []*Worksheet

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Teams []*Team

type Photo struct {
	ID            int
	WorksheetID   int
	RunningNumber int
	FileName      string
	Location      string
	Created       time.Time
}

type Photos []*Photo

func (f *Photo) FilePath() string {
	return "/store/" + strconv.Itoa(f.WorksheetID) + "/" + f.FileName
}

type FormField struct {
	ID    int
	Name  string
	Title string
}

type FormFields []*FormField
