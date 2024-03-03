import { configureStore } from "@reduxjs/toolkit";
import authReducer from "@/slices/authSlice";

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
