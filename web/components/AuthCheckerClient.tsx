"use client";

import React, { useMemo, useState } from 'react';
import { TokenData, TokenDataElId, UserData, UserDataElId } from '@/lib/auth/Auth';

import { AuthSlice } from '@/slices/authSlice';
import { LocalApi } from '@/lib/api/LocalApi';
import authCookie from '@/lib/auth/AuthCookie';
import { getCookies } from 'cookies-next';
import { useAppDispatch } from "@/store/hooks";

const getTokenData = (): TokenData | undefined => {
    if (typeof document === "undefined")
        return undefined;

    const tokenEl = document.getElementById(TokenDataElId) as HTMLInputElement;
    if (tokenEl) {
        const tokenJson = Buffer.from(tokenEl.value, 'hex').toString();
        const tokenData = JSON.parse(tokenJson) as TokenData;
        if (tokenData.id_token) {
            console.log("got token data from hidden el", tokenData);
            // tokenEl.remove();
            return tokenData;
        }
    }
}

const getUserData = (): UserData | undefined => {
    if (typeof document === "undefined")
        return undefined;

    const userEl = document.getElementById(UserDataElId) as HTMLInputElement;
    if (userEl) {
        const userJson = Buffer.from(userEl.value, 'hex').toString();
        const userData = JSON.parse(userJson) as UserData;
        if (userData.user_id) {
            console.log("got user data from hidden el", userData);
            // userEl.remove();
            return userData;
        }
    }
}

export const AuthCheckerClient = () => {   
    const dispatch = useAppDispatch();
    const [called, setCalled] = useState(false);

    // Redefined only to prevent confusion with useMemo
    const useInit = (callback: () => unknown, depends = []) => useMemo(callback, depends) 

    useInit(() => {
        if (!called) {

            const tokenData = getTokenData();
            const userData = getUserData();

            if (userData) {
                // dispatch to save user state
                dispatch(AuthSlice.actions.setUserData(userData));
            }

            if (tokenData) {
                dispatch(AuthSlice.actions.setAuthToken(tokenData));
                // LocalApi.getInstance().setAuthCookie(tokenData).then();
            }

            setCalled(true);
        }
    }, []);

    return (
        <></>
    );
}
