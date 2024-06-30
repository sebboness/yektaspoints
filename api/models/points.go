package models

import (
	"slices"
	"time"

	"github.com/sebboness/yektaspoints/util"
)

type PointRequestType string

const PointRequestTypeAdd PointRequestType = "ADD"
const PointRequestTypeSubtract PointRequestType = "SUBTRACT"
const PointRequestTypeCashout PointRequestType = "CASHOUT"

var PointSubtractTypes = []PointRequestType{
	PointRequestTypeSubtract,
	PointRequestTypeCashout,
}

type PointRequestDecision string

const PointRequestDecisionApprove PointRequestDecision = "APPROVE"
const PointRequestDecisionDeny PointRequestDecision = "DENY"

type PointStatus string

const PointStatusWaiting = "WAITING"
const PointStatusSettled = "SETTLED"

var ValidPointRequestDecisions = []PointRequestDecision{
	PointRequestDecisionApprove,
	PointRequestDecisionDeny,
}

type Point struct {
	ID           string       `json:"id" dynamodbav:"id"`
	UserID       string       `json:"user_id" dynamodbav:"user_id"`
	Status       PointStatus  `json:"status" dynamodbav:"status"`
	Points       int32        `json:"points" dynamodbav:"points"`
	Balance      *int32       `json:"balance" dynamodbav:"balance,omitempty"`
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
	ID              string               `json:"id"`
	UserID          string               `json:"user_id"`
	ParentNotes     string               `json:"parent_notes"`
	Reason          string               `json:"reason"`
	Points          int32                `json:"points"`
	Type            PointRequestType     `json:"type"`
	UpdatedOn       time.Time            `json:"updated_on"`
	DecidedByUserID string               `json:"decided_by_user_id"`
	Decision        PointRequestDecision `json:"decision" dynamodbav:"decision,omitempty"`
}

type UserPoints struct {
	Balance             int32          `json:"balance"`
	PointsLast7Days     int32          `json:"points_last_7_days"`
	PointsLostLast7Days int32          `json:"points_lost_last_7_days"`
	RecentCashouts      []PointSummary `json:"recent_cashouts"`
	RecentRequests      []PointSummary `json:"recent_requests"`
	RecentPoints        []PointSummary `json:"recent_points"`
}

type PointBalance struct {
	ID      string `json:"id" dynamodbav:"id"`
	UserID  string `json:"user_id" dynamodbav:"user_id"`
	Balance int32  `json:"balance" dynamodbav:"balance"`
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

func (p *Point) ToPointSummary() PointSummary {
	return PointSummary{
		ID:              p.ID,
		UserID:          p.UserID,
		ParentNotes:     p.Request.ParentNotes,
		Points:          p.Points,
		Reason:          p.Request.Reason,
		UpdatedOn:       p.UpdatedOn,
		Type:            p.Request.Type,
		DecidedByUserID: p.Request.DecidedByUserID,
		Decision:        p.Request.Decision,
	}
}

func ToPointSummaries(points []Point) []PointSummary {
	summaries := make([]PointSummary, len(points))
	for idx, p := range points {
		summaries[idx] = p.ToPointSummary()
	}
	return summaries
}

// IsSubtractType returns true if the given point request type is a subtraction type (i.e. SUBTRACT or CASHOUT)
func IsSubtractType(reqType PointRequestType) bool {
	return slices.Contains(PointSubtractTypes, reqType)
}
