import { Point, UserPoints } from "@/lib/models/Points";
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";

import { ErrorAsResult } from "@/lib/api/Result";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { getTokenRetriever } from "@/store/store";
import moment from "moment";

const logName = () => `[${moment().toISOString()}] pointsSlice: `;

type LoginOptions = {
    username: string;
    password: string;
}

export const getUserPointSummary = createAsyncThunk("auth/login", async (userID: string, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api
            .withToken(getTokenRetriever())
            .getPointSummaryByUser(userID);

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
        // setAuthToken: (state, action: PayloadAction<TokenData | undefined>) => {
        //     console.log(`${logName()}setAuthToken: token`, action.payload);
        //     state.token = action.payload;
        // },
    },
    extraReducers: (builder) => {
        builder.addCase(getUserPointSummary.fulfilled, (state, action) => {
            state.userSummary = action.payload;
        });
        builder.addCase(getUserPointSummary.rejected, (state, action) => {
            console.log(`${logName()}getUserPointSummary rejected`, action.error);
            // TODO hmmm
        });
    },
});

export default PointsSlice.reducer;
