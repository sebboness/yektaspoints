package models

type Family struct {
	FamilyID string                  `json:"family_id"`
	Children map[string]FamilyMember `json:"children"`
	Parents  map[string]FamilyMember `json:"parents"`
}

type FamilyUser struct {
	FamilyID string `json:"family_id" dynamodbav:"family_id"`
	UserID   string `json:"user_id" dynamodbav:"user_id"`
}

type FamilyMember struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`

	// Name used in app and displayed to children (i.e. "Mom")
	ChildCallName string `json:"child_call_name"`
}

func NewFamilyUser(user User) FamilyMember {
	return FamilyMember{
		Email:         user.Email,
		UserID:        user.UserID,
		Name:          user.Name,
		ChildCallName: user.ChildCallName,
	}
}
