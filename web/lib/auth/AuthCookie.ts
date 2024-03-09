import { RequestCookies, ResponseCookies } from "next/dist/compiled/@edge-runtime/cookies";
import { TokenData, UserData } from "./Auth";

import { NextResponse } from "next/server";
import { ReadonlyRequestCookies } from "next/dist/server/web/spec-extension/adapters/request-cookies";
import { cookies } from "next/headers";
import moment from "moment";

const ln = () => `[${moment().toISOString()}] AuthCookie: `;

export const TokenCookieName = "mypoints_web_auth";
export const RefreshCookieName = "mypoints_web_rt";
export const UserCookieName = "mypoints_web_usr";

class AuthCookie {
    constructor() {
    }

    deleteAll(response: NextResponse, domain: string): NextResponse {
        const env = process.env["ENV"];

        // cookie settings
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`${ln()}deleteAll ${TokenCookieName} cookie on ${env}:${_domain}`);
        console.info(`${ln()}deleteAll ${RefreshCookieName} cookie on ${env}:${_domain}`);
        console.info(`${ln()}deleteAll ${UserCookieName} cookie on ${env}:${_domain}`);

        // Unset token data cookie
        response.cookies.set({
            name: TokenCookieName,
            value: "",
            maxAge: -1, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Unset refresh token cookie
        response.cookies.set({
            name: RefreshCookieName,
            value: "",
            maxAge: -1, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Unset user cookie
        response.cookies.set({
            name: UserCookieName,
            value: "",
            maxAge: -1, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        return response;
    }

    async set(domain: string, tokenData: TokenData): Promise<unknown> {
        const refreshToken = tokenData.refresh_token;
        tokenData.refresh_token = "";
        const tokenJson = JSON.stringify(tokenData);
        const env = process.env["ENV"];

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`${ln()}set ${TokenCookieName} cookie on ${env}:${_domain}`); // with value ${tokenJson}`);
        console.info(`${ln()}set ${RefreshCookieName} cookie on ${env}:${_domain}`); // with value ${refreshToken}`);

        cookies().set({
            name: TokenCookieName,
            value: tokenJson,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Set refresh token cookie
        cookies().set({
            name: RefreshCookieName,
            value: refreshToken,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        return;
    }

    async setTokenData(cookies: ResponseCookies | RequestCookies, domain: string, tokenData: TokenData) {
        const refreshToken = tokenData.refresh_token;
        tokenData.refresh_token = "";
        const tokenJson = JSON.stringify(tokenData);
        const env = process.env["ENV"];

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`${ln()}setTokenData ${TokenCookieName} cookie on ${env}:${_domain}`); // with value ${tokenJson}`);
        console.info(`${ln()}setTokenData ${RefreshCookieName} cookie on ${env}:${_domain}`); // with value ${refreshToken}`);

        // Set token data cookie (without refresh token value)
        cookies.set({
            name: TokenCookieName,
            value: tokenJson,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Set refresh token cookie
        cookies.set({
            name: RefreshCookieName,
            value: refreshToken,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });
    }

    async setUserData(cookies: ResponseCookies | RequestCookies, domain: string, userData: UserData) {
        const userJson = JSON.stringify(userData);
        const env = process.env["ENV"];

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`${ln()}setUserData ${UserCookieName} cookie on ${env}:${_domain} with value ${userJson}`);

        // Set user data cookie
        cookies.set({
            name: UserCookieName,
            value: userJson,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
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
