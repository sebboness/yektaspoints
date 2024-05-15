import React from "react";
import { cookies } from "next/headers";

import { TokenDataElId, UserDataElId } from "@/lib/auth/Auth";
import authCookie from "@/lib/auth/AuthCookie";

export const AuthChecker = async () => {

    const tokenData = authCookie.getTokenData(cookies());
    const userData = authCookie.getUserData(cookies());

    const tokenHexed = Buffer.from(JSON.stringify(tokenData || {})).toString("hex");
    const userHexed = Buffer.from(JSON.stringify(userData || {})).toString("hex");

    return (
        <>
            <><input type="hidden" id={TokenDataElId} value={tokenHexed} /></>
            <><input type="hidden" id={UserDataElId} value={userHexed} /></>
        </>
    );
};
