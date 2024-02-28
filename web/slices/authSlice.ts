import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";

export const fetchUsers = createAsyncThunk("auth/fetchUsers", async (thunkApi) => {
    console.log("getting users...");
    const response = await fetch("https://jsonplaceholder.typicode.com/users?_limit=50");
    const data = await response.json();
    console.log("got users.");
    return data;
})

type UserData = {
    roles: string[];
    userId: string;
    username: string;
    email: string;
    name: string;
}

type AuthState = {
    accessToken: string;
    idToken: string;
    expiresIn: number;
    user: UserData | null;
}

const initialState: AuthState = {
    accessToken: "",
    idToken: "",
    expiresIn: 0,
    user: null,
};

const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        increment: (state) => {
            state.accessToken = "123";
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
    },
});

export default authSlice.reducer;
