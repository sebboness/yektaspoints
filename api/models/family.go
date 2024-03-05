package models

type Family struct {
	FamilyID string       `json:"family_id"`
	Children []FamilyUser `json:"children"`
	Parents  []FamilyUser `json:"parents"`
}

type FamilyUser struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`

	// Name used in app and displayed to children (i.e. "Mom")
	ChildCallName string `json:"child_call_name"`
}
