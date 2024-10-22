package forms

import "strings"

type LoginUser struct {
	Username string
	Password string
	Failures map[string]string
}

func (f *LoginUser) Valid() bool {
	f.Failures = make(map[string]string)
	if strings.TrimSpace(f.Username) == "" {
		f.Failures["Username"] = "Username is required"
	}
	if strings.TrimSpace(f.Password) == "" {
		f.Failures["Password"] = "Password is required"
	}
	return len(f.Failures) == 0
}

type Query struct {
	Q            string
	Date         string
	TicketTypeID int
	Start        int
	MaxResults   int
}

func NewQuery() *Query {
	return &Query{MaxResults: 100}
}

type Team struct {
	Name string `form:"team_name"`
}

type Zone struct {
	Name string `form:"zone_name"`
}

type Worksheet struct {
	Number string `form:"worksheet_number"`
	Name   string `form:"worksheet_name"`
	ZoneID int    `form:"worksheet_zone_id"`
	TeamID int    `form:"worksheet_team_id"`
}

type File struct {
	RunningNumber string `form:"running_number"`
}

type Attendee struct {
	Firstname    string            `form:"firstname"`
	Lastname     string            `form:"lastname"`
	TicketTypeID int               `form:"ticket_type_id"`
	CustomFields map[string]string `form:"custom_fields"`
}
