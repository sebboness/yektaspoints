package models

import "time"

type DateFilter struct {
	From *time.Time
	To   *time.Time
}

// IsSet returns true if either from or to dates are not nil (not set); False otherwise
func (f DateFilter) IsSet() bool {
	return !(f.From == nil && f.To == nil)
}
