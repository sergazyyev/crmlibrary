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

//All ad_groups
const (
	//Block
	AdGroupBsBmb = "CRM_BS_BMB"
	AdGroupBsBkb = "CRM_BS_BKB"
	AdGroupBsBrb = "CRM_BS_BRB"
	AdGroupBsBsb = "CRM_BS_BSB"

	//Channel
	AdGroupChPole       = "CRM_CH_POLE"
	AdGroupChOnline     = "CRM_CH_ONLINE"
	AdGroupChAgent      = "CRM_CH_AGENT"
	AdGroupChTm         = "CRM_CH_TM"
	AdGroupChController = "CRM_CH_CONTROLLER"
	AdGroupChCallCenter = "CRM_CH_CALLCENTER"
	AdGroupChUpm        = "CRM_CH_UPM"

	//Role
	AdGroupRlEmployee = "CRM_RL_EMPLOYEE"
	AdGroupRlAnalytic = "CRM_RL_ANALYTIC"
	AdGroupRlManager  = "CRM_RL_MANAGER"
	AdGroupRlChief    = "CRM_RL_CHIEF"
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
