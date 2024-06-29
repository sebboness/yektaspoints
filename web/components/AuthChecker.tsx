import React from "react";
import moment from "moment";
import { cookies, headers } from "next/headers";

import { TokenDataElId, UserDataElId } from "@/lib/auth/Auth";
import authCookie from "@/lib/auth/AuthCookie";

const ln = () => `[${moment().toISOString()}] AuthChecker: `;

export const AuthChecker = async () => {

    const tokenData = authCookie.getTokenDataFromHeader(headers());
    const userData = authCookie.getUserData(cookies());

    console.log(`${ln()}token? ${tokenData ? (tokenData.id_token.substring(tokenData.id_token.length - 20)) : "NONE"}`);

    const tokenHexed = Buffer.from(JSON.stringify(tokenData || {})).toString("hex");
    const userHexed = Buffer.from(JSON.stringify(userData || {})).toString("hex");

    return (
        <>
            <><input type="hidden" id={TokenDataElId} value={tokenHexed} /></>
            <><input type="hidden" id={UserDataElId} value={userHexed} /></>
        </>
    );
};
