import { RequestCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { NextResponse } from "next/server";
import { ReadonlyRequestCookies } from "next/dist/server/web/spec-extension/adapters/request-cookies";
import moment from "moment";

import { TokenData, UserData } from "./Auth";

const ln = () => `[${moment().toISOString()}] AuthCookie: `;

export const TokenCookieName = "mypoints_web_auth";
export const RefreshCookieName = "mypoints_web_rt";
export const UserCookieName = "mypoints_web_usr";
export const TokenHeaderName = "X-Points4Us-Api-Token";

class AuthCookie {
    constructor() {
    }

    deleteAll(res: NextResponse, domain: string): NextResponse {
        const env = process.env.ENV;

        // cookie settings
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`${ln()}deleteAll ${TokenCookieName} cookie on ${env}:${_domain}`);
        console.info(`${ln()}deleteAll ${RefreshCookieName} cookie on ${env}:${_domain}`);
        console.info(`${ln()}deleteAll ${UserCookieName} cookie on ${env}:${_domain}`);

        // Unset token data cookie
        res.cookies.set({
            name: TokenCookieName,
            value: "",
            maxAge: -1, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
            path: "/",
        });

        // Unset refresh token cookie
        res.cookies.set({
            name: RefreshCookieName,
            value: "",
            maxAge: -1, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
            path: "/",
        });

        // Unset user cookie
        res.cookies.set({
            name: UserCookieName,
            value: "",
            maxAge: -1, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
            path: "/",
        });

        return res;
    }

    setTokenData(res: NextResponse, domain: string, tokenData: TokenData) {
        const refreshToken = tokenData.refresh_token;
        tokenData.refresh_token = "";
        const tokenJson = JSON.stringify(tokenData);
        const env = process.env.ENV;

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`${ln()}setTokenData ${TokenCookieName} cookie on ${env}:${_domain}`); // with value ${tokenJson}`);
        console.info(`${ln()}setTokenData ${RefreshCookieName} cookie on ${env}:${_domain}`); // with value ${refreshToken}`);

        // Set token data cookie (without refresh token value)
        res.cookies.set({
            name: TokenCookieName,
            value: tokenJson,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
            path: "/",
        });

        // Set refresh token cookie
        res.cookies.set({
            name: RefreshCookieName,
            value: refreshToken,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
            path: "/",
        });
    }

    setTokenDataToHeader(headers: Headers, tokenData: TokenData) {
        // remove unused access_token key
        delete tokenData["access_token"];
        const tokenJson = JSON.stringify(tokenData);
        const env = process.env.ENV;

        console.info(`${ln()}setTokenDataToHeader ${TokenHeaderName} header on ${env}`); // with value ${tokenJson}`);

        // Set token data header
        headers.set(TokenHeaderName, tokenJson);
    }

    setUserData(res: NextResponse, domain: string, userData: UserData) {
        const userJson = JSON.stringify(userData);
        const env = process.env.ENV;

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`${ln()}setUserData ${UserCookieName} cookie on ${env}:${_domain} with value ${userJson}`);

        // Set user data cookie
        res.cookies.set({
            name: UserCookieName,
            value: userJson,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
            path: "/",
        });
    }

    getTokenData(cookies: ReadonlyRequestCookies | RequestCookies): TokenData | undefined {
        const tokenCookie = cookies.get(TokenCookieName);
        const refreshTokenCookie = cookies.get(RefreshCookieName);
        
        console.info(`${ln()}cookie ${TokenCookieName} set? ${tokenCookie !== undefined}`);
        console.info(`${ln()}cookie ${RefreshCookieName} set? ${refreshTokenCookie !== undefined}`);

        if (!tokenCookie || !refreshTokenCookie)
            return undefined;

        const tokenJson = tokenCookie.value;
        const tokenData = JSON.parse(tokenJson) as TokenData;
        if (!tokenData.id_token) {
            console.info("token cookie not a valid TokenData object: " + tokenJson);
            return undefined;
        }

        tokenData.refresh_token = refreshTokenCookie.value;

        return tokenData;
    }

    getTokenDataFromHeader(headers: Headers): TokenData | undefined {
        const tokenJson = headers.get(TokenHeaderName);
        
        console.info(`${ln()}header ${TokenCookieName} set? ${tokenJson !== undefined}`);

        if (!tokenJson)
            return undefined;

        const tokenData = JSON.parse(tokenJson) as TokenData;
        if (!tokenData.id_token) {
            console.info("token header not a valid TokenData object: " + tokenJson);
            return undefined;
        }

        return tokenData;
    }

    getUserData(cookies: ReadonlyRequestCookies | RequestCookies): UserData | undefined {
        const userCookie = cookies.get(UserCookieName);
        
        console.info(`${ln()}cookie ${UserCookieName} set? ${userCookie !== undefined}`);

        if (!userCookie)
            return undefined;

        const userJson = userCookie.value;
        const userData = JSON.parse(userJson) as UserData;
        if (!userData.user_id) {
            console.info("user cookie not a valid UserData object: " + userJson);
            return undefined;
        }

        return userData;
    }
}

const authCookie = new AuthCookie();

export default authCookie;
