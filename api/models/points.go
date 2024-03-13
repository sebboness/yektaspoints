package models

import (
	"time"

	"github.com/sebboness/yektaspoints/util"
)

type PointRequestType string

const PointRequestTypeAdd PointRequestType = "ADD"
const PointRequestTypeSubtract PointRequestType = "SUBTRACT"
const PointRequestTypeCashout PointRequestType = "CASHOUT"

type PointRequestDecision string

const PointRequestDecisionApprove PointRequestDecision = "APPROVE"
const PointRequestDecisionDeny PointRequestDecision = "DENY"

type PointStatus string

const PointStatusWaiting = "WAITING"
const PointStatusSettled = "SETTLED"

type Point struct {
	ID           string       `json:"id" dynamodbav:"id"`
	UserID       string       `json:"user_id" dynamodbav:"user_id"`
	Status       PointStatus  `json:"status" dynamodbav:"status"`
	Points       int          `json:"points" dynamodbav:"points"`
	Balance      *int         `json:"balance" dynamodbav:"balance,omitempty"`
	CreatedOnStr string       `json:"-" dynamodbav:"created_on"`
	UpdatedOnStr string       `json:"-" dynamodbav:"updated_on"`
	CreatedOn    time.Time    `json:"created_on" dynamodbav:"-"`
	UpdatedOn    time.Time    `json:"updated_on" dynamodbav:"-"`
	Request      PointRequest `json:"request" dynamodbav:"request"`
}

type PointRequest struct {
	DecidedByUserID string               `json:"decided_by_user_id" dynamodbav:"decided_by_user_id,omitempty"`
	DecidedOnStr    string               `json:"-" dynamodbav:"decided_on,omitempty"`
	DecidedOn       time.Time            `json:"decided_on" dynamodbav:"-"`
	Decision        PointRequestDecision `json:"decision" dynamodbav:"decision,omitempty"`
	ParentNotes     string               `json:"parent_notes" dynamodbav:"parent_notes,omitempty"`
	Reason          string               `json:"reason" dynamodbav:"reason,omitempty"`
	Type            PointRequestType     `json:"type" dynamodbav:"type"`
}

type QueryPointsFilter struct {
	CreatedOn  DateFilter
	UpdatedOn  DateFilter
	Statuses   []PointStatus
	Types      []PointRequestType
	Attributes []string // Which attributes to project in the query
}

type PointSummary struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	ParentNotes     string           `json:"parent_notes"`
	Reason          string           `json:"reason"`
	UpdatedOn       time.Time        `json:"updated_on"`
	Type            PointRequestType `json:"type"`
	DecidedByUserID string           `json:"decided_by_user_id"`
}

type UserPoints struct {
	Balance        int            `json:"balance"`
	RecentCashouts []PointSummary `json:"recent_cashouts"`
	RecentRequests []PointSummary `json:"recent_requests"`
	RecentPoints   []PointSummary `json:"recent_points"`
}

func (p *Point) ParseTimes() {
	if p.CreatedOnStr != "" {
		p.CreatedOn = util.ParseTime_RFC3339Nano(p.CreatedOnStr)
	}
	if p.UpdatedOnStr != "" {
		p.UpdatedOn = util.ParseTime_RFC3339Nano(p.UpdatedOnStr)
	}
	if p.Request.DecidedOnStr != "" {
		p.Request.DecidedOn = util.ParseTime_RFC3339Nano(p.Request.DecidedOnStr)
	}
}
