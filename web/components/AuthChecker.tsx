import { MyPointsApi } from '@/lib/api/MyPointsApi';
import { FAILURE } from '@/lib/api/Result';
import { ParseToken } from '@/lib/auth/Auth';
import authCookie from '@/lib/auth/AuthCookie';
import { getSimpleTokenRetriever } from '@/slices/authSlice';
import React from 'react'
import { redirect } from 'next/navigation';
import { cookies, headers } from 'next/headers';
import moment from 'moment';
import { Roles } from '@/lib/auth/Roles';
import _ from 'lodash';
import { encodeHtmlEntities } from '@/lib/StringHelpers';

const logName = () => `[${moment().toISOString()}] AuthChecker: `;

type Props = {
    roles?: Roles[];
}

export const AuthChecker = async (props: Props) => {

    // console.log("cookies", cookies().get("mypoints_web_auth"));
    // const cookieToken = getTokenFromCookie();

    let idToken = "";
    let rtToken = "";
    let username = "";
    let roles: Roles[] = [];

    async function setCookieServerSide() {
        "use server"

        cookies().set({
            name: "hey",
            value: moment().toISOString(),
            maxAge: 1000, 
            httpOnly: true,
            sameSite: "strict",
            secure: false,
            domain: "localhost",
        });
    };

    const tokenData = await authCookie.get();
    if (!tokenData) {
        console.warn(`${logName()}no token cookie. issue redirect to login`);
        redirect("/login");
    } else {
        console.info(`${logName()}tokenCookie has data`);

        const userData = ParseToken(tokenData.id_token);
        if (!userData) {
            console.warn(`${logName()}failed to parse token to user data.`);
            redirect("/access-denied?user");
        }

        roles = userData.groups;
        username = userData.username;
        idToken = tokenData.id_token;
        rtToken = tokenData.refresh_token;
        
        // now check against mypoints api if cookie is valid
        const authResp = await MyPointsApi.getInstance()
            .withToken(getSimpleTokenRetriever(tokenData.id_token))
            .getUserAuth();

        if (authResp.status === FAILURE) {
            console.info(`${logName()}getUserAuth response`, JSON.stringify(authResp));
            console.info(`${logName()}getUserAuth failed with ${userData.username} ${tokenData.id_token}`);
            console.info(`${logName()}trying to refresh with ${userData.username} ${tokenData.refresh_token}`);

            // token has possibly expired, so try to refresh it.
            const refreshResp = await MyPointsApi.getInstance()
                .refreshToken(userData.username, tokenData.refresh_token);

            if (refreshResp.data) {
                idToken = refreshResp.data.id_token;
                rtToken = refreshResp.data.refresh_token;
            } else {
                console.info(`${logName()}refreshToken response`, JSON.stringify(refreshResp));
                console.info(`${logName()}need to redirect to login page`);

                redirect("/login");
            }
        }
    }

    if (props.roles && !_.find(props.roles, (x) => _.find(roles, (y) => y === x))) {
        console.info(`${logName()}required roles [${props.roles}] not found in user roles [${roles}]`);
        redirect("/access-denied?roles");
    }

    // const store = useAppStore();
    // const userLoggedIn = false;
    // const authState = useAppSelector((state) => state.auth);

    // // console.log("state", store.getState());
    // console.log("authState", authState);

    // // is user defined? if not, we know user is not logged in
    // if (!authState.user) {
    //     console.info("no auth data in state");
    //     // const localApi = LocalApi.getInstance().getAuthCookie()
    //     //     .then(())
    //     const tokenData = getTokenFromCookie();
    //     if (tokenData) {
    //         console.info("got token data from cookie");
    //         // checks if token in cookie is still valid
    //         // store.dispatch(checkUserAuth(tokenData));
    //     } else {
    //         console.info("no token data in cookie");
    //         redirect("/login");
    //     }
    // } else {
    //     // console.info("auth data in state", store.getState().auth);
    // }

    const b = Buffer.from(JSON.stringify(tokenData)).toString('hex');

    return (
        <><pre style={{"display": "none"}} id="__tokendata__">{b}</pre></>
    )
}
