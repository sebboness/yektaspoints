import { TokenGetter } from "@/lib/api/Api";
import { LocalApi } from "@/lib/api/LocalApi";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { ErrorAsResult, SUCCESS } from "@/lib/api/Result";
import { ParseToken, TokenData, UserData } from "@/lib/auth/Auth";
import { createSlice, createAsyncThunk, PayloadAction } from "@reduxjs/toolkit";

/**
 * Checks if given token data is still valid
 */
export const checkUserAuth = createAsyncThunk("auth/checkUserAuth", async (tokenData: TokenData, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api.withToken(getSimpleTokenRetriever(tokenData.id_token)).getUserAuth();
        if (result.status === SUCCESS) {
            thunkApi.dispatch(authSlice.actions.setAuthToken(tokenData));
            return tokenData;
        }
        else
            throw result; // this means we should try to refresh
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

/**
 * Clears auth cookie
 */
const clearAuthCookie = createAsyncThunk("auth/clearAuthCookie", async (params, thunkApi) => {
    const api = LocalApi.getInstance();
    const resp = await api.deleteAuthCookie();
    return resp.status === SUCCESS;
})

/**
 * Sets auth cookie
 */
const setAuthCookie = createAsyncThunk("auth/setAuthCookie", async (tokenData: TokenData, thunkApi) => {
    const api = LocalApi.getInstance();
    const resp = await api.setAuthCookie(tokenData);
    return resp.status === SUCCESS;
})

/**
 * Gets a token getter that returns the token that was passed into it.
 * @returns The token getter
 */
export const getSimpleTokenRetriever = (token: string): TokenGetter => {
    return {
        getToken() {return token; },
        getTokenType()  { return "Bearer"; },
    };
}

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
        if (result.status === SUCCESS) {
            thunkApi.dispatch(authSlice.actions.setAuthToken(result.data!));
            return result.data!;
        }
        else
            throw result;
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

type AuthState = {
    token?: TokenData;
    user?: UserData;
    authCookieSet: boolean;
};

const initialState: AuthState = {
    authCookieSet: false,
};

export const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        setAuthToken: (state, action: PayloadAction<TokenData>) => {
            console.log("authSlice.setAuthToken: setting token");
            const user = ParseToken(action.payload.id_token);
            console.log("authSlice.setAuthToken: user", user);
            state.token = action.payload;
            state.user = user;
        },
    },
    extraReducers: (builder) => {
        builder.addCase(checkUserAuth.rejected, (state, action) => {
            console.log("checkUserAuth rejected", action.error);

            state.token = undefined;
            state.user = undefined;
        });

        builder.addCase(login.rejected, (state, action) => {
            console.log("login rejected", action.error);

            state.token = undefined;
            state.user = undefined;
        });
    },
});

export default authSlice.reducer;
