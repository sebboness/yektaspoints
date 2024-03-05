import { MyPointsApi } from '@/lib/api/MyPointsApi';
import { FAILURE } from '@/lib/api/Result';
import { ParseToken } from '@/lib/auth/Auth';
import authCookie from '@/lib/auth/AuthCookie';
import { getSimpleTokenRetriever } from '@/slices/authSlice';
import React from 'react'
import { LocalApi } from '@/lib/api/LocalApi';
import { redirect } from 'next/navigation';

const logName = () => "AuthChecker: ";

export const AuthChecker = async () => {

    // console.log("cookies", cookies().get("mypoints_web_auth"));
    // const cookieToken = getTokenFromCookie();

    let idToken = "";
    let rtToken = "";
    let username = "";

    const tokenData = await authCookie.get();
    if (!tokenData) {
        console.warn(`${logName()}no token cookie. issue redirect to login`);
        redirect("/login");
    } else {
        console.info(`${logName()}tokenCookie has data`);
        idToken = tokenData.id_token;
        rtToken = tokenData.refresh_token;
        
        // now check against mypoints api if cookie is valid
        const authResp = await MyPointsApi.getInstance()
            .withToken(getSimpleTokenRetriever(tokenData.id_token))
            .getUserAuth();

        if (authResp.status === FAILURE) {
            console.info(`${logName()}getUserAuth response`, JSON.stringify(authResp));

            // token has possibly expired, so try to refresh it.
            const userData = ParseToken(tokenData.id_token);
            const refreshResp = await MyPointsApi.getInstance()
                .refreshToken(userData?.username!, tokenData.refresh_token);

            username = userData?.username!;

            if (refreshResp.data) {
                console.info(`${logName()}successfully refreshed. Setting token cookie...`);
                await LocalApi.getInstance().setAuthCookie(refreshResp.data);

                idToken = refreshResp.data.id_token;
                rtToken = refreshResp.data.refresh_token;
            } else {
                console.info(`${logName()}refreshToken response`, JSON.stringify(refreshResp));
                console.info(`${logName()}need to redirect to login page`);

                redirect("/login");
            }
        }

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

    return (
        <></>
    )
}
