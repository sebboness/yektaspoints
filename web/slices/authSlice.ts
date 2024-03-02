import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { ErrorAsResult, SUCCESS } from "@/lib/api/Result";
import { TokenData, UserData } from "@/lib/auth/Auth";
import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";

export const fetchUsers = createAsyncThunk("auth/fetchUsers", async (params, thunkApi) => {
    // thunkApi.dispatch(login({username: "", password: ""}));
    console.log("getting users...");
    const response = await fetch("https://jsonplaceholder.typicode.com/users?_limit=50");
    const data = await response.json();
    console.log("got users.");
    return data;
});

type LoginOptions = {
    username: string;
    password: string;
}

type RefreshOptions = {
    username: string;
    refresh_token: string;
}

export const login = createAsyncThunk("auth/login", async (options: LoginOptions, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api.authenticate(options.username, options.password);
        if (result.status === SUCCESS)
            return result.data!;
        else
            throw result;
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

type AuthState = {
    token?: TokenData;
    user?: UserData;
};

const initialState: AuthState = {
};

const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        increment: (state) => {
            // state.accessToken = "123";
        },
    },
    extraReducers: (builder) => {
        builder.addCase(fetchUsers.fulfilled, (state, action) => {
            // state.entities.push(...action.payload);
            // state.loading = false;
        });

        builder.addCase(fetchUsers.pending, (state, action) => {
            // state.loading = true;
        });

        builder.addCase(login.fulfilled, (state, action) => {
            console.log("login fulfilled", action.payload);
            state.token = action.payload

        });

        builder.addCase(login.rejected, (state, action) => {
            console.log("login rejected", action.error);
        });
    },
});

export default authSlice.reducer;
