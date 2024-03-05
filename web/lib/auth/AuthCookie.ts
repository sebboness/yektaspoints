import { cookies } from "next/headers";
import { TokenData } from "./Auth";
import { compress, decompress } from "shrink-string";
import { NextResponse } from "next/server";

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
        const rtCompressed = await compress(tokenData.refresh_token);
        tokenData.refresh_token = "";
        const tokenJson = JSON.stringify(tokenData);
        const tokenCompressed = await compress(tokenJson);
        const env = process.env["ENV"];

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`Setting ${tokenCookieName} cookie on ${env}:${_domain} with value ${tokenCompressed}`);
        console.info(`Setting ${refreshCookieName} cookie on ${env}:${_domain} with value ${rtCompressed}`);

        cookies().set({
            name: tokenCookieName,
            value: tokenCompressed,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Set refresh token cookie
        cookies().set({
            name: refreshCookieName,
            value: rtCompressed,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        return;
    }

    async setWithResponse(response: NextResponse, domain: string, tokenData: TokenData): Promise<NextResponse> {
        const rtCompressed = await compress(tokenData.refresh_token);
        tokenData.refresh_token = "";
        const tokenJson = JSON.stringify(tokenData);
        const tokenCompressed = await compress(tokenJson);
        const env = process.env["ENV"];

        // cookie settings
        const maxAge = 60*60*24*30;// 30 days
        const secure = env == "local" ? false : true;
        const _domain = domain === "localhost" ? domain : "hexonite.net";

        console.info(`Setting ${tokenCookieName} cookie on ${env}:${_domain} with value ${tokenCompressed}`);
        console.info(`Setting ${refreshCookieName} cookie on ${env}:${_domain} with value ${rtCompressed}`);

        // Set token data cookie (without refresh token value)
        response.cookies.set({
            name: tokenCookieName,
            value: tokenCompressed,
            maxAge: maxAge, 
            httpOnly: true,
            sameSite: "strict",
            secure: secure,
            domain: _domain,
        });

        // Set refresh token cookie
        response.cookies.set({
            name: refreshCookieName,
            value: rtCompressed,
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

        const tokenJson = await decompress(tokenCookie.value);
        const refreshToken = await decompress(refreshTokenCookie.value);

        const tokenData = JSON.parse(tokenJson) as TokenData;
        tokenData.refresh_token = refreshToken;

        return tokenData;
    }
}

const authCookie = new AuthCookie();

export default authCookie;
