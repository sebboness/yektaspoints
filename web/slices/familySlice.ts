import moment from "moment";

import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";

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
        setFamilies: (state, action: PayloadAction<MapType<Family>>) => {
            console.log(`${ln()}setFamilies`, action.payload);
            state.families = action.payload;
        },
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
