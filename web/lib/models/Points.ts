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
