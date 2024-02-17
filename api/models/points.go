package models

import (
	"time"

	"github.com/sebboness/yektaspoints/util"
)

type PointType string

const PointTypeAdd PointType = "ADD"
const PointTypeSubtract PointType = "SUBTRACT"
const PointTypeCashout PointType = "CASHOUT"
const PointTypeWallet PointType = "WALLET"

type PointStatus int

const PointStatusRequested = 0
const PointStatusApproved = 1
const PointStatusDenied = -1

type Point struct {
	ID              string      `json:"id" dynamodbav:"id"`
	UserID          string      `json:"userId" dynamodbav:"user_id"`
	DecidedByUserID string      `json:"decidedByUserId" dynamodbav:"decided_by_user_id,omitempty"`
	StatusID        PointStatus `json:"statusId" dynamodbav:"status_id"`
	Points          int         `json:"points" dynamodbav:"points"`
	BalancePoints   int         `json:"balance_points" dynamodbav:"balance_points"`
	Balance         int         `json:"balance" dynamodbav:"balance"`
	Reason          string      `json:"reason" dynamodbav:"reason,omitempty"`
	Notes           string      `json:"notes" dynamodbav:"notes,omitempty"`
	Type            PointType   `json:"type" dynamodbav:"type"`
	RequestedOnStr  string      `json:"-" dynamodbav:"requested_on"`
	DecidedOnStr    string      `json:"-" dynamodbav:"decided_on,omitempty"`
	RequestedOn     time.Time   `json:"requestedOn" dynamodbav:"-"`
	DecidedOn       time.Time   `json:"decidedOn" dynamodbav:"-"`
}

type QueryPointsFilter struct {
	RequestedOn DateFilter
	Statuses    []PointStatus
	Types       []PointType
}

func (p Point) ParseTimes() {
	if p.RequestedOnStr != "" {
		p.RequestedOn = util.ParseTime_RFC3339Nano(p.RequestedOnStr)
	}
	if p.DecidedOnStr != "" {
		p.DecidedOn = util.ParseTime_RFC3339Nano(p.DecidedOnStr)
	}
}
