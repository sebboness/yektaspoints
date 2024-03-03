import { configureStore } from "@reduxjs/toolkit";
import authReducer from "@/slices/authSlice";
import { TokenGetter } from "@/lib/api/Api";

export function makeStore() {
    return configureStore({
        reducer: {
            auth: authReducer,
        },
    })
}

export const store = makeStore();

export type AppStore = ReturnType<typeof makeStore>
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

/**
 * Gets a token getter that returns the auth token from the auth store.
 * @returns The token getter
 */
export const getTokenRetriever = (): TokenGetter => {
    return {
        getToken() {
            const authStore = store.getState().auth;
            return authStore.token
                ? authStore.token.id_token
                : "";
        },

        getTokenType()  { return "Bearer"},
    };
}