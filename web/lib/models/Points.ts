export enum PointDecision {
    APPROVE = "APPROVE",
    DENY = "DENY",
};

export enum PointStatus {
    SETTLED = "SETTLED",
    WAITING = "WAITING",
};

export enum PointRequestType {
    ADD = "ADD",
    SUBTRACT = "SUBTRACT",
    CASHOUT = "CASHOUT",
};

export type Point = {
    id: string
    user_id: string
    status: string
    points: number
    balance?: number
    created_on: string
    updated_on: string
    request: PointRequest
}

export type PointsList = {
    points: Point[];
}

export interface PointRequest {
    decided_by_user_id?: string
    decided_on?: string
    decision?: string
    parent_notes?: string
    reason?: string
    type: string
}

export type PointSummary = {
    id: string
    parent_notes?: string
    reason?: string
    points: number
    type: string
    updated_on: string
    decided_by_user_id?: string
    decision?: string
};

export type UserPoints = {
    balance: number
    points_last_7_days: number
    points_lost_last_7_days: number
    recent_cashouts: PointSummary[]
    recent_requests: PointSummary[]
    recent_points: PointSummary[]
};

export type RequestPointsRequest = {
    points: number;
    reason: string;
}

export type RequestPointsResponse = {
    point: Point;
    point_summary: PointSummary;
}

export const mapPointToSummary = (p: Point): PointSummary => ({
    id: p.id,
    parent_notes: p.request.parent_notes,
    reason: p.request.reason,
    points: p.points,
    type: p.request.type,
    updated_on: p.updated_on,
    decided_by_user_id: p.request.decided_by_user_id,
    decision: p.request.decision,
});

export const mapPointsToSummaries = (points: Point[]): PointSummary[] => 
    Object.values(points).map((p => mapPointToSummary(p)));
