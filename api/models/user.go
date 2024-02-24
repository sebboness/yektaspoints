package models

import "time"

type User struct {
	UserID    int       `json:"user_id" dynamodbav:"user_id"`
	Name      string    `json:"name" dynamodbav:"name"`
	Username  string    `json:"username" dynamodbav:"username"`
	Email     string    `json:"email" dynamodbav:"email"`
	ChildIds  []string  `json:"child_ids" dynamodbav:"child_ids"`
	ParentIds []string  `json:"parent_ids" dynamodbav:"parent_ids"`
	CreatedOn time.Time `json:"created_on" dynamodbav:"created_on"`
	UpdatedOn time.Time `json:"updated_on" dynamodbav:"updated_on"`
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
