package storage

type UserDTO struct {
	ID        int    `db:"id"`
	StatusID  int    `db:"statusId"`
	Name      string `db:"name"`
	Username  string `db:"username"`
	Password  string `db:"password"`
	CreatedOn string `db:"createdOn"`
	UpdatedOn string `db:"updatedOn"`
}
