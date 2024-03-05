import { cookies } from "next/headers";
import { NextResponse } from "next/server";

import { TokenData } from "./Auth";

const tokenCookieName = "mypoints_web_auth";
const refreshCookieName = "mypoints_web_rt";

class AuthCookie {
    constructor() {
    }

    delete(response: NextResponse, domain: string): NextResponse {
        const env = process.env["ENV"];

        // cookie settings
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`Deleting ${tokenCookieName} cookie on ${env}:${_domain}`);
        console.info(`Deleting ${refreshCookieName} cookie on ${env}:${_domain}`);

        // Set token data cookie (without refresh token value)
        response.cookies.set({
            name: tokenCookieName,
            value: "",
            maxAge: -1, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Set refresh token cookie
        response.cookies.set({
            name: refreshCookieName,
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

        console.info(`Setting ${tokenCookieName} cookie on ${env}:${_domain} with value ${tokenJson}`);
        console.info(`Setting ${refreshCookieName} cookie on ${env}:${_domain} with value ${refreshToken}`);

        cookies().set({
            name: tokenCookieName,
            value: tokenJson,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Set refresh token cookie
        cookies().set({
            name: refreshCookieName,
            value: refreshToken,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        return;
    }

    async setWithResponse(response: NextResponse, domain: string, tokenData: TokenData): Promise<NextResponse> {
        const refreshToken = tokenData.refresh_token;
        tokenData.refresh_token = "";
        const tokenJson = JSON.stringify(tokenData);
        const env = process.env["ENV"];

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`Setting ${tokenCookieName} cookie on ${env}:${_domain} with value ${tokenJson}`);
        console.info(`Setting ${refreshCookieName} cookie on ${env}:${_domain} with value ${refreshToken}`);

        // Set token data cookie (without refresh token value)
        response.cookies.set({
            name: tokenCookieName,
            value: tokenJson,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Set refresh token cookie
        response.cookies.set({
            name: refreshCookieName,
            value: refreshToken,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        return response;
    }

    async get(): Promise<TokenData | undefined> {
        const tokenCookie = cookies().get(tokenCookieName);
        const refreshTokenCookie = cookies().get(refreshCookieName);
        
        console.info(`cookie ${tokenCookieName} set? ${tokenCookie !== undefined}`);
        console.info(`cookie ${refreshCookieName} set? ${refreshTokenCookie !== undefined}`);

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
}

const authCookie = new AuthCookie();

export default authCookie;
