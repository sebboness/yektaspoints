import { MyPointsApi } from '@/lib/api/MyPointsApi';
import authCookie from '@/lib/auth/AuthCookie';
import { getSimpleTokenRetriever } from '@/slices/authSlice';
import React from 'react'

export const AuthChecker = async () => {

    // console.log("cookies", cookies().get("mypoints_web_auth"));
    // const cookieToken = getTokenFromCookie();

    const tokenData = await authCookie.get();
    console.info("tokenCookie", tokenData);
    if (!tokenData) {
        console.warn("no token cookie. issue redirect to login");
    } else {
        console.info("token data from cookie:", JSON.stringify(tokenData));
        const authResp = await MyPointsApi.getInstance()
            .withToken(getSimpleTokenRetriever(tokenData.id_token))
            .getUserAuth();

        console.info("check user auth response:", JSON.stringify(authResp));

        // now check against mypoints api if cookie is valid
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
