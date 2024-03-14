package models

import "time"

type DateFilter struct {
	From *time.Time
	To   *time.Time
}

func NewDateFilter() *DateFilter {
	return &DateFilter{}
}

func (df *DateFilter) WithRange(from time.Time, to time.Time) *DateFilter {
	df.From = &from
	df.To = &to
	return df
}

func (df *DateFilter) WithFrom(t time.Time) *DateFilter {
	df.From = &t
	return df
}

func (df *DateFilter) WithTo(t time.Time) *DateFilter {
	df.To = &t
	return df
}

// IsSet returns true if either from or to dates are not nil (not set); False otherwise
func (f DateFilter) IsSet() bool {
	return !(f.From == nil && f.To == nil)
}
