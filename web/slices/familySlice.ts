import moment from "moment";

import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";

import { Family } from "@/lib/models/Family";
import { ErrorAsResult } from "@/lib/api/Result";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { MapType } from "@/lib/models/Common";

const ln = () => `[${moment().toISOString()}] familySlice: `;

export const getFamily = createAsyncThunk("getFamily", async (familyId: string, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api.getFamily(familyId);

        if (result.data)
            return result.data;
        else
            throw result;
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

type FamilyState = {
    families: MapType<Family>;
};

const initialState: FamilyState = {
    families: {},
};

export const FamilySlice = createSlice({
    name: "family",
    initialState,
    reducers: {
        // addPointToRequesting: (state, action: PayloadAction<PointSummary>) => {
        //     state.userSummary.recent_requests.unshift(action.payload);
        //     console.log(`${ln()}addPointToRequesting payload`, action.payload);
        //     console.log(`${ln()}addPointToRequesting new recent point reqs`, state.userSummary.recent_requests);
        // },
    },
    extraReducers: (builder) => {
        builder.addCase(getFamily.fulfilled, (state, action) => {
            state.families[action.payload.family.family_id] = action.payload.family;
        });
        builder.addCase(getFamily.rejected, (state, action) => {
            console.log(`${ln()}getFamily rejected`, action.error);
            // TODO hmmm
        });
    },
});

export default FamilySlice.reducer;
