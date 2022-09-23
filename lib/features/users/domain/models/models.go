package models

import "time"

type DetailedUser struct {
	Id          string    `json:"id"`
	FullName    string    `json:"full_name"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	BirthDate   time.Time `json:"birth_date"` // always 0 minutes and 0 seconds
	Address     Address   `json:"address"`
	// Preferences Preferences
}

type Address struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Nation  string `json:"nation"`
	ZipCode string `json:"zipcode"`
}

// type Preferences struct {
// 	Language string
//   DateFormat dateFormatPref
// }
//
// type dateFormatPref int
// const (
// 	MonthFirstFormat dateFormatPref = iota
// 	DayFirstFormat
// )
