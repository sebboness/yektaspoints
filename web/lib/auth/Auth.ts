import { jwtDecode } from "jwt-decode";
import { getCookie } from "cookies-next";

export const TokenCookieName = "mypoints_web_auth";

export type TokenData = {
    access_token: string;
    id_token: string;
    refresh_token: string;
    expires_in: number;
    new_password_required: boolean;
    session: string;
}

export type UserData = {
    email: string;
    name: string;
    groups: string[];
    user_id: string;
    username: string;
    verified: boolean;
    exp: number;
};

type JwtData = {[key: string]: any};

/**
 * Parses the given token and returns user data
 * @param token The JWT token to decode
 * @returns Parse user data from decoded token
 */
export const ParseToken = (token: string): UserData | undefined => {
    try {
        const data = jwtDecode(token) as JwtData;
        return {
            email: data["email"],
            name: data["name"],
            groups: data["cognito:groups"] || [],
            user_id: data["sub"],
            username: data["cognito:username"],
            verified: data["email_verified"],
            exp: data["exp"],
        };
    } catch (err) {
        console.warn(`failed to decode jwt token: ${token}`);
    }

    return undefined;
}

export const getTokenFromCookie = (): TokenData | undefined => {
    const tokenCookie = getCookie(TokenCookieName);
    if (!tokenCookie)
        return undefined;

    return JSON.parse(tokenCookie) as TokenData;
}