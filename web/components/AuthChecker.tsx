import { TokenDataElId, UserDataElId } from "@/lib/auth/Auth";

import React from "react";
import { Roles } from "@/lib/auth/Roles";
import _ from "lodash";
import authCookie from "@/lib/auth/AuthCookie";
import { cookies } from "next/headers";
import moment from "moment";

const logName = () => `[${moment().toISOString()}] AuthChecker: `;

type Props = {
    roles?: Roles[];
};

export const AuthChecker = async (props: Props) => {

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
