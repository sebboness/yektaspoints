import { NextRequest, NextResponse } from "next/server";
import { ParseToken, TokenData, UserData } from "@/lib/auth/Auth";

import { FAILURE } from "@/lib/api/Result";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { Roles } from "@/lib/auth/Roles";
import authCookie, { TokenHeaderName } from "@/lib/auth/AuthCookie";
import { cookies } from "next/headers";
import { getSimpleTokenRetriever } from "@/slices/authSlice";
import moment from "moment";
import { xRedirectToHeader } from "./config";

const ln = () => `[${moment().toISOString()}] AuthMiddleware: `;
const parentsBasePath = ["/family", "/app/parents"];
const childrenBasePath = ["/points", "/app/mypoints"];
const adminBasePath = ["/admin"];

export default async function AuthMiddleware(req: NextRequest): Promise<NextResponse> {

    const path = req.nextUrl.pathname;
    
    if (
        !parentsBasePath.some((p) => path.startsWith(p)) &&
        !childrenBasePath.some((p) => path.startsWith(p)) &&
        !adminBasePath.some((p) => path.startsWith(p))
    ) {
        // No match on logged-in paths, so pass it on
        // console.log(`${ln()}no match on path ${path}: ${req.nextUrl.href}`);
        return NextResponse.next();
    } else {
        console.log(`${ln()}path match on ${path}: ${req.nextUrl.href}`);
        let userData: UserData | undefined;
        let tokenData: TokenData | undefined;
        let isLoginRedirect = false;

        // Check cookies
        tokenData = authCookie.getTokenData(cookies());
        if (!tokenData) {
            // No auth cookie set, so redirect to login page
            console.warn(`${ln()}no token cookie. issue redirect to login`);
            isLoginRedirect = true;
        } else {
            // Cookie present, so lets parse the token and grab the user data
            console.log(`${ln()}tokenCookie has data`);

            userData = ParseToken(tokenData.id_token);
            if (!userData) {
                // Couldn't parse token from cookie, so redirect to login page
                console.warn(`${ln()}failed to parse token user data, so redirect to login`);
                isLoginRedirect = true;
            } else {
                // User data is parsed, and now let's check token against auth API
                const username = userData.username;

                let getUserResp = await MyPointsApi.getInstance()
                    .withToken(getSimpleTokenRetriever(tokenData.id_token))
                    .getUser();

                if (getUserResp.status === FAILURE || !getUserResp.data) {
                    // Token has possibly expired, so try to refresh it.
                    console.log(`${ln()}getUser response`, JSON.stringify(getUserResp));
                    console.log(`${ln()}getUser failed for ${username}`); // with token ${tokenData.id_token}`);
                    console.log(`${ln()}trying to refresh for ${username}`); // with refresh token ${tokenData.refresh_token}`);

                    // invalidate user data
                    userData = undefined;

                    const refreshResp = await MyPointsApi.getInstance()
                        .refreshToken(username, tokenData.refresh_token);

                    if (!refreshResp.data) {
                        // Refresh token failed, so redirect to login page
                        console.log(`${ln()}refreshToken response`);
                        console.log(`${ln()}need to redirect to login page`);
                        isLoginRedirect = true;

                        // invalidate token
                        tokenData = undefined;
                    } else {
                        // Now that we refreshed the auth token, let's get user data again
                        getUserResp = await MyPointsApi.getInstance()
                            .withToken(getSimpleTokenRetriever(refreshResp.data.id_token))
                            .getUser();

                        if (!getUserResp.data) {
                            // This shouldn't happen, but getting authed user data failed, so redirect to login page
                            // and log an error with this
                            console.warn(`${ln()}getUser response`, JSON.stringify(getUserResp));
                            console.warn(`${ln()}need to redirect to login page`);
                            isLoginRedirect = true;
                        } else {
                            // All good!
                            console.log(`${ln()}token and user data all good!`);
                            userData = getUserResp.data;
                            tokenData = refreshResp.data;
                        }
                    }
                } else {
                    userData = getUserResp.data;
                }
            }
        }

        const response = NextResponse.next();
        const roles: Roles[] = userData
            ? (userData.roles || [])
            : [];

        // Set token data cookie
        if (tokenData) {
            console.log(`${ln()}setting token? ${tokenData.id_token.substring(tokenData.id_token.length - 20)}`);
            console.log(`${ln()}A`);
            authCookie.setTokenData(response, req.nextUrl.hostname, tokenData);
            response.headers.set(TokenHeaderName, tokenData.id_token);
        }

        // Set user data cookie and check
        if (userData) {
            authCookie.setUserData(response, req.nextUrl.hostname, userData);
        }

        // Check first if we should issue a login redirect, then check specific user roles for requested path
        if (isLoginRedirect) {
            response.headers.set(xRedirectToHeader, `/login?return_url=${encodeURIComponent(path)}`);
        } else if (parentsBasePath.some((p) => path.startsWith(p)) && !roles.some((r) => r === Roles.Parent)) {
            // No parent role
            console.log(`${ln()}trying to access parents page, but is missing role [${roles}]`);
            response.headers.set(xRedirectToHeader, "/access-denied");
        } else if (childrenBasePath.some((p) => path.startsWith(p)) && !roles.some((r) => r === Roles.Child)) {
            // No child role
            console.log(`${ln()}trying to access kids page, but is missing role [${roles}]`);
            response.headers.set(xRedirectToHeader, "/access-denied");
        } else if (adminBasePath.some((p) => path.startsWith(p)) && !roles.some((r) => r === Roles.Admin)) {
            // No child role
            console.log(`${ln()}trying to access admin page, but is missing role [${roles}]`);
            response.headers.set(xRedirectToHeader, "/access-denied");
        }
        
        console.log(`${ln()}${response.headers.has(xRedirectToHeader) ? ("redirect to " + response.headers.get(xRedirectToHeader)) : "no redirects"}`);        

        return response;
    }
}