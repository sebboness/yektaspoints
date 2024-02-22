package models

import "time"

type UserRegisterResponse struct {
	ID        int       `json:"id"`
	StatusID  int       `json:"statusId"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"createdOn"`
	UpdatedOn time.Time `json:"updatedOn"`
}
