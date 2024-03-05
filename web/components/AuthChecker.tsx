import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { FAILURE } from "@/lib/api/Result";
import { ParseToken } from "@/lib/auth/Auth";
import authCookie from "@/lib/auth/AuthCookie";
import { Roles } from "@/lib/auth/Roles";
import { getSimpleTokenRetriever } from "@/slices/authSlice";
import _ from "lodash";
import moment from "moment";
import { redirect } from "next/navigation";
import React from "react";

const logName = () => `[${moment().toISOString()}] AuthChecker: `;

type Props = {
    roles?: Roles[];
};

export const AuthChecker = async (props: Props) => {

    let roles: Roles[] = [];

    const tokenData = await authCookie.get();
    if (!tokenData) {
        console.warn(`${logName()}no token cookie. issue redirect to login`);
        redirect("/login");
    } else {
        console.log(`${logName()}tokenCookie has data`);

        const userData = ParseToken(tokenData.id_token);
        if (!userData) {
            console.warn(`${logName()}failed to parse token to user data.`);
            redirect("/access-denied?user");
        }

        roles = userData.groups;

        // now check against mypoints api if cookie is valid
        const authResp = await MyPointsApi.getInstance()
            .withToken(getSimpleTokenRetriever(tokenData.id_token))
            .getUserAuth();

        if (authResp.status === FAILURE) {
            console.log(`${logName()}getUserAuth response`, JSON.stringify(authResp));
            console.log(`${logName()}getUserAuth failed with ${userData.username} ${tokenData.id_token}`);
            console.log(`${logName()}trying to refresh with ${userData.username} ${tokenData.refresh_token}`);

            // token has possibly expired, so try to refresh it.
            const refreshResp = await MyPointsApi.getInstance()
                .refreshToken(userData.username, tokenData.refresh_token);

            if (!refreshResp.data) {
                console.log(`${logName()}refreshToken response`, JSON.stringify(refreshResp));
                console.log(`${logName()}need to redirect to login page`);

                redirect("/login");
            }
        }
    }

    if (props.roles && !_.find(props.roles, (x) => _.find(roles, (y) => y === x))) {
        console.log(`${logName()}required roles [${props.roles}] not found in user roles [${roles}]`);
        redirect("/access-denied?roles");
    }

    const tokenHexed = Buffer.from(JSON.stringify(tokenData)).toString("hex");

    return (
        <><pre style={{"display": "none"}} id="__tokendata__">{tokenHexed}</pre></>
    );
};
