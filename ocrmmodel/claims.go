package ocrmmodel

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sergazyyev/crmlibrary/ocrmerrors"
)

const (
	GroupTypeBlock   = "BLOCK"
	GroupTypeRole    = "ROLE"
	GroupTypeChannel = "CHANNEL"
)

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

func (claims *Claims) GetUserGroupsByClaims() (block string, role string, channel string, err error) {
	groups := claims.User.Groups
	for _, group := range groups {
		switch group.Type {
		case GroupTypeBlock:
			block = group.Name
		case GroupTypeChannel:
			channel = group.Name
		case GroupTypeRole:
			role = group.Name
		}
	}
	if block == "" || role == "" || channel == "" {
		err = ocrmerrors.New(ocrmerrors.INVALID, "Can't parse users roles by token", "Невозможно распарсить роли пользователя по токену")
		return
	}
	return
}
