package ocrmmodel

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	User *User `json:"user"`
	*jwt.StandardClaims
}

type User struct {
	Username   string    `json:"username"`
	Mail       string    `json:"mail"`
	Name       string    `json:"name"`
	Department string    `json:"department"`
	Manager    string    `json:"manager"`
	JobTitle   string    `json:"jobTitle"`
	City       string    `json:"city"`
	Groups     []*Group  `json:"groups"`
	Modules    []*Module `json:"modules"`
}

type Group struct {
	Name  string `json:"name"`
	Type  string `db:"type" json:"type"`
	Level int    `db:"level" json:"level"`
}

type Module struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}
