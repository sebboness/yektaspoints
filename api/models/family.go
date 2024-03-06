package models

type Family struct {
	FamilyID string                `json:"family_id"`
	Children map[string]FamilyUser `json:"children"`
	Parents  map[string]FamilyUser `json:"parents"`
}

type FamilyUser struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`

	// Name used in app and displayed to children (i.e. "Mom")
	ChildCallName string `json:"child_call_name"`
}

func NewFamilyUser(user User) FamilyUser {
	return FamilyUser{
		Email:         user.Email,
		UserID:        user.UserID,
		Name:          user.Name,
		ChildCallName: user.ChildCallName,
	}
}
