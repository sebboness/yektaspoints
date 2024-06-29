import { AuthCookieBody, ParseToken, TokenData, UserData } from "@/lib/auth/Auth";
import { ErrorAsResult, SUCCESS } from "@/lib/api/Result";
import { PayloadAction, createAsyncThunk, createSlice } from "@reduxjs/toolkit";

import { RootState } from "@/store/store";
import { LocalApi } from "@/lib/api/LocalApi";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { TokenGetter } from "@/lib/api/Api";
import moment from "moment";

const ln = () => `[${moment().toISOString()}] authSlice: `;

let refreshTimer: NodeJS.Timeout | undefined;
const refreshBuffer = 3 * 60 * 1000; // refresh in 3 minutes before the actual time out in milliseconds

/**
 * Clears auth cookie
 */
const clearAuthCookie = createAsyncThunk("auth/clearAuthCookie", async (params, thunkApi) => {
    try {
        const api = LocalApi.getInstance();
        const resp = await api.deleteAuthCookie();
        return resp.status === SUCCESS;
    }
    catch (err) {
        // log a warning but return true nonetheless to clear data in slice
        console.warn(`${ln()}error clearing auth cookie`);
        return true;
    }
});

export const startRefreshTimer = createAsyncThunk("auth/startRefreshTimer", async (expiresAt: number, thunkApi) => {
    if (expiresAt <= 0) {
        console.log(`${ln()}refresh timer expiresAt is not set`);
        return;
    }

    console.log(`${ln()}starting refresh timer. overwriting existing timer?`, refreshTimer !== undefined);

    // clear original timout first
    if (refreshTimer)
        clearTimeout(refreshTimer);

    const now = new Date().getTime();
    const refreshInMs = (expiresAt * 1000) - now - refreshBuffer;
    console.log(`${ln()}${expiresAt} - ${now} - ${refreshBuffer} = ${refreshInMs}`);

    if (refreshInMs <= 0) {
        console.log(`${ln()}refreshInMsis too little? ${expiresAt} - ${now} - ${refreshBuffer} = ${refreshInMs}`);
        return;
    }

    refreshTimer = setTimeout(async () => {
        const state = thunkApi.getState() as RootState;
        const user = state.auth.user;
        const token = state.auth.token;
        console.log(`${ln()}dispatching user refresh`, user && token);
        if (user && token) {
            await thunkApi.dispatch(refresh({
                username: user.username,
                refreshToken: token.refresh_token,
            }));
        }
    }, refreshInMs);
});

/**
 * Gets user data from api with currently logged in auth token
 */
export const getUser = createAsyncThunk("getUser", async (params, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api.getUser();

        if (result.data)
            return result.data;
        else
            throw result;
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

/**
 * Sets auth cookie for token and user
 */
export const setAuthCookie = createAsyncThunk("auth/setAuthCookie", async (authCookie: AuthCookieBody, thunkApi) => {
    const api = LocalApi.getInstance();
    const resp = await api.setAuthCookie(authCookie);
    return resp.status === SUCCESS;
});

/**
 * Gets a token getter that returns the token that was passed into it.
 * @returns The token getter
 */
export const getSimpleTokenRetriever = (token: string): TokenGetter => {
    return {
        getToken() { return token; },
        getTokenType() { return "Bearer"; },
    };
};

type LoginOptions = {
    username: string;
    password: string;
};

type RefreshOptions = {
    username: string;
    refreshToken: string;
};

export const login = createAsyncThunk("auth/login", async (options: LoginOptions, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const result = await api.authenticate(options.username, options.password);
        if (result.status === SUCCESS) {
            thunkApi.dispatch(AuthSlice.actions.setAuthToken(result.data!));
            return result.data!;
        }
        else
            thunkApi.dispatch(clearAuthCookie());
    } catch (err: any) {
        throw ErrorAsResult(err);
    }
});

export const refresh = createAsyncThunk("auth/refresh", async (options: RefreshOptions, thunkApi) => {
    const api = MyPointsApi.getInstance();
    try {
        const refreshResult = await api.refreshToken(options.username, options.refreshToken);
        if (refreshResult.status === SUCCESS && refreshResult.data) {
            const newToken = ParseToken(refreshResult.data.id_token || "");
            const expiresAt = newToken ? newToken.exp : 0;

            const userResult = await MyPointsApi.getInstance()
                .withToken(getSimpleTokenRetriever(refreshResult.data.id_token))
                .getUser();

            if (userResult.status === SUCCESS && userResult.data) {
                userResult.data.exp = expiresAt;

                thunkApi.dispatch(setAuthCookie({
                    token: refreshResult.data,
                    user: userResult.data,
                }));

                thunkApi.dispatch(AuthSlice.actions.setAuthToken(refreshResult.data));
                thunkApi.dispatch(AuthSlice.actions.setUserData(userResult.data));
                thunkApi.dispatch(startRefreshTimer(expiresAt));

                return refreshResult.data!;
            }
            else
                thunkApi.dispatch(clearAuthCookie());
        }
        else
            thunkApi.dispatch(clearAuthCookie());
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

export const AuthSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        setAuthToken: (state, action: PayloadAction<TokenData | undefined>) => {
            console.log(`${ln()}setAuthToken: token`, action.payload);
            state.token = action.payload;
            MyPointsApi.getInstance()
                .withToken(getSimpleTokenRetriever(action.payload?.id_token || ""));
        },

        setUserData: (state, action: PayloadAction<UserData | undefined>) => {
            console.log(`${ln()}setUserData: user`, action.payload);
            state.user = action.payload;
        },
    },
    extraReducers: (builder) => {
        builder.addCase(clearAuthCookie.fulfilled, (state, action) => {
            console.log(`${ln()}clearAuthCookie fulfilled`);
            state.token = undefined;
            state.user = undefined;
        });

        builder.addCase(login.rejected, (state, action) => {
            console.log(`${ln()}login rejected`, action.error);
            state.token = undefined;
            state.user = undefined;
        });

        builder.addCase(refresh.rejected, (state, action) => {
            console.log(`${ln()}refresh rejected`, action.error);
            state.token = undefined;
            state.user = undefined;
        });

        builder.addCase(setAuthCookie.fulfilled, (state, action) => {
            state.authCookieSet = action.payload;
        });
        builder.addCase(setAuthCookie.rejected, (state, action) => {
            console.log(`${ln()}setAuthCookie rejected`, action.error);
            state.authCookieSet = false;
        });
        
        builder.addCase(getUser.fulfilled, (state, action) => {
            state.user = action.payload;
        });
    },
});

export default AuthSlice.reducer;
