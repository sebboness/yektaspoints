"use client";

import { decodeHtmlEntities } from '@/lib/StringHelpers';
import { LocalApi } from '@/lib/api/LocalApi';
import { TokenData } from '@/lib/auth/Auth';
import React from 'react';

export const AuthCheckerClient = () => {

    if (typeof window !== 'undefined') {
        const tokenEl = document.getElementById("__tokendata__");
        if (tokenEl) {
            const tokenJson = Buffer.from(tokenEl.innerHTML, 'hex').toString();
            const tokenData = JSON.parse(tokenJson) as TokenData;
            if (tokenData.id_token) {
                console.log("yoohoo we found tokendata!", tokenData);
                
                LocalApi.getInstance().setAuthCookie(tokenData)
                    .then((result) => console.info("token data sent to be saved in cookie", result))
                    .catch((err) => console.warn("failed to store token data cookie on backend", err));
            }
        }
    }

    return (
        <div>AuthCheckerClient</div>
    )
}
