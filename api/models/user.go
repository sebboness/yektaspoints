package models

import (
	"slices"
	"time"

	"github.com/sebboness/yektaspoints/util"
)

type UserStatus string

const UserStatusActive UserStatus = "ACTIVE"
const UserStatusDeleted UserStatus = "DELETED"
const UserStatusInactive UserStatus = "INACTIVE"
const UserStatusUnverified UserStatus = "UNVERIFIED"

type User struct {
	Email        string     `json:"email" dynamodbav:"email"`
	FamilyIDs    []string   `json:"family_ids" dynamodbav:"family_ids"`
	Name         string     `json:"name" dynamodbav:"name"`
	UserID       string     `json:"user_id" dynamodbav:"user_id"`
	Username     string     `json:"username" dynamodbav:"username"`
	Roles        []string   `json:"roles" dynamodbav:"roles"`
	Status       UserStatus `json:"status" dynamodbav:"status"`
	CreatedOnStr string     `json:"-" dynamodbav:"created_on"`
	UpdatedOnStr string     `json:"-" dynamodbav:"updated_on,omitempty"`
	CreatedOn    time.Time  `json:"created_on" dynamodbav:"-"`
	UpdatedOn    time.Time  `json:"updated_on" dynamodbav:"-"`

	// Name used in app and displayed to children (i.e. "Mom")
	ChildCallName string `json:"child_call_name" dynamodbav:"child_call_name"`
}

func (u *User) ParseTimes() {
	if u.CreatedOnStr != "" {
		u.CreatedOn = util.ParseTime_RFC3339Nano(u.CreatedOnStr)
	}
	if u.UpdatedOnStr != "" {
		u.UpdatedOn = util.ParseTime_RFC3339Nano(u.UpdatedOnStr)
	}
}

func (u *User) IsAdmin() bool {
	return slices.Contains(u.Roles, "admin")
}

func (u *User) IsChild() bool {
	return slices.Contains(u.Roles, "child")
}

func (u *User) IsParent() bool {
	return slices.Contains(u.Roles, "parent")
}
