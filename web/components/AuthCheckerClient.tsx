"use client";

import { LocalApi } from '@/lib/api/LocalApi';
import { TokenData } from '@/lib/auth/Auth';
import React, { useEffect } from 'react';

export const AuthCheckerClient = () => {    

    useEffect(() => {
        if (typeof window !== 'undefined') {
            const tokenEl = document.getElementById("__tokendata__");
            if (tokenEl) {
                const tokenJson = Buffer.from(tokenEl.innerHTML, 'hex').toString();
                const tokenData = JSON.parse(tokenJson) as TokenData;
                if (tokenData.id_token) {                
                    LocalApi.getInstance().setAuthCookie(tokenData)
                        .then((result) => console.info("token data sent to be saved in cookie", result))
                        .catch((err) => console.warn("failed to store token data cookie on backend", err));

                    // remove token data from the dom
                    tokenEl.remove();
                }
            }
        }
    }, []);

    return (
        <div>AuthCheckerClient</div>
    )
}
