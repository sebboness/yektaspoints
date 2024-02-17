package models

import "time"

type DateFilter struct {
	From *time.Time
	To   *time.Time
}

// IsSet returns true if both from and to dates are nil (not set); True otherwise
func (f DateFilter) IsSet() bool {
	return !(f.From == nil && f.To == nil)
}
