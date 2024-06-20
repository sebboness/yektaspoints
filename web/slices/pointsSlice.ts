import { PayloadAction, createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { Point, PointSummary, UserPoints } from "@/lib/models/Points";

import { getTokenRetriever } from "@/store/store";
import { ErrorAsResult } from "@/lib/api/Result";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import moment from "moment";

const ln = () => `[${moment().toISOString()}] pointsSlice: `;

export const getUserPoints = createAsyncThunk("getUserPoints", async (userID: string, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api
            .withToken(getTokenRetriever())
            .getUserPoints(userID);

        if (result.data)
            return result.data;
        else
            throw result;
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

export const getUserPointSummary = createAsyncThunk("getUserPointSummary", async (userID: string, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api
            .withToken(getTokenRetriever())
            .getUserPointsSummary(userID);

        if (result.data)
            return result.data;
        else
            throw result;
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

type PointsState = {
    userPoints: Point[];
    userSummary: UserPoints;
};

const initialState: PointsState = {
    userPoints: [],
    userSummary: {
        balance: 0,
        points_last_7_days: 0,
        points_lost_last_7_days: 0,
        recent_cashouts: [],
        recent_points: [],
        recent_requests: [],
    }
};

export const PointsSlice = createSlice({
    name: "points",
    initialState,
    reducers: {
        addPointToRequesting: (state, action: PayloadAction<PointSummary>) => {
            state.userSummary.recent_requests.unshift(action.payload);
            console.log(`${ln()}addPointToRequesting payload`, action.payload);
            console.log(`${ln()}addPointToRequesting new recent point reqs`, state.userSummary.recent_requests);
        },
    },
    extraReducers: (builder) => {
        builder.addCase(getUserPoints.fulfilled, (state, action) => {
            state.userPoints = action.payload.points;
        });
        builder.addCase(getUserPoints.rejected, (state, action) => {
            console.log(`${ln()}getUserPoints rejected`, action.error);
            // TODO hmmm
        });
        builder.addCase(getUserPointSummary.fulfilled, (state, action) => {
            state.userSummary = action.payload;
        });
        builder.addCase(getUserPointSummary.rejected, (state, action) => {
            console.log(`${ln()}getUserPointSummary rejected`, action.error);
            // TODO hmmm
        });
    },
});

export default PointsSlice.reducer;
