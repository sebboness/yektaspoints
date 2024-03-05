package models

import "time"

type User struct {
	Email     string    `json:"email" dynamodbav:"email"`
	FamilyID  string    `json:"family_id" dynamodbav:"family_id"`
	UserID    string    `json:"user_id" dynamodbav:"user_id"`
	Username  string    `json:"username" dynamodbav:"username"`
	Name      string    `json:"name" dynamodbav:"name"`
	Roles     []string  `json:"roles" dynamodbav:"roles"`
	CreatedOn time.Time `json:"created_on" dynamodbav:"created_on"`
	UpdatedOn time.Time `json:"updated_on" dynamodbav:"updated_on"`

	// Name used in app and displayed to children (i.e. "Mom")
	ChildCallName string `json:"child_call_name" dynamodbav:"child_call_name"`
}

type UserRegisterResponse struct {
	ID        int       `json:"id"`
	StatusID  int       `json:"statusId"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"createdOn"`
	UpdatedOn time.Time `json:"updatedOn"`
}
