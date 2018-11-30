package models

import (
	"strconv"
	"time"
)

type Worksheet struct {
	ID       int    `json:"id"`
	Number   string `json:"number"`
	Name     string `json:"name"`
	ZoneID   int
	ZoneName string
	TeamID   int
	TeamName string
	Created  time.Time `json:"created"`
}

type Worksheets []*Worksheet

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Teams []*Team

type Zone struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Zones []*Zone

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

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Locations []*Location
