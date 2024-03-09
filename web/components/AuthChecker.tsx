import { ParseToken, TokenDataElId, UserData, UserDataElId } from "@/lib/auth/Auth";

import { FAILURE } from "@/lib/api/Result";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import React from "react";
import { Roles } from "@/lib/auth/Roles";
import _ from "lodash";
import authCookie from "@/lib/auth/AuthCookie";
import { cookies } from "next/headers";
import { getSimpleTokenRetriever } from "@/slices/authSlice";
import moment from "moment";
import { redirect } from "next/navigation";

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

    // let roles: Roles[] = [];
    // let userData: UserData | undefined = undefined;

    // const tokenData = await authCookie.getTokenData(cookies());
    // if (!tokenData) {
    //     console.warn(`${logName()}no token cookie. issue redirect to login`);
    //     redirect("/login");
    // } else {
    //     console.log(`${logName()}tokenCookie has data`);

    //     userData = ParseToken(tokenData.id_token);
    //     if (!userData) {
    //         console.warn(`${logName()}failed to parse token to user data.`);
    //         redirect("/access-denied?user");
    //     }

    //     roles = userData.roles;

    //     // now check against mypoints api if cookie is valid
    //     let authResp = await MyPointsApi.getInstance()
    //         .withToken(getSimpleTokenRetriever(tokenData.id_token))
    //         .getUser();

    //     if (authResp.status === FAILURE || !authResp.data) {
    //         console.log(`${logName()}getUserAuth response`, JSON.stringify(authResp));
    //         console.log(`${logName()}getUserAuth failed with ${userData.username} ${tokenData.id_token}`);
    //         console.log(`${logName()}trying to refresh with ${userData.username} ${tokenData.refresh_token}`);

    //         // token has possibly expired, so try to refresh it.
    //         const refreshResp = await MyPointsApi.getInstance()
    //             .refreshToken(userData.username, tokenData.refresh_token);

    //         if (!refreshResp.data) {
    //             console.log(`${logName()}refreshToken response`, JSON.stringify(refreshResp));
    //             console.log(`${logName()}need to redirect to login page`);
    //             redirect("/login");
    //         }

    //         // Now that we refreshed, let's get user data again
    //         authResp = await MyPointsApi.getInstance()
    //             .withToken(getSimpleTokenRetriever(refreshResp.data.id_token))
    //             .getUser();

    //         if (!authResp.data) {
    //             // failed again, so redirect to login
    //             console.log(`${logName()}authResp response`, JSON.stringify(refreshResp));
    //             console.log(`${logName()}need to redirect to login page`);
    //             redirect("/login");
    //         }
    //     }

    //     userData = authResp.data;
    // }

    // if (props.roles && !_.find(props.roles, (x) => _.find(roles, (y) => y === x))) {
    //     console.log(`${logName()}required roles [${props.roles}] not found in user roles [${roles}]`);
    //     redirect("/access-denied?roles");
    // }

    // const tokenHexed = Buffer.from(JSON.stringify(tokenData)).toString("hex");
    // const userHexed = Buffer.from(JSON.stringify(userData)).toString("hex");

    // return (
    //     <>
    //         <><input type="hidden" id={TokenDataElId} value={tokenHexed} /></>
    //         <><input type="hidden" id={UserDataElId} value={userHexed} /></>
    //     </>
    // );
};
