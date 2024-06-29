import { jwtDecode } from "jwt-decode";

import { Roles } from "./Roles";

export const TokenDataElId = "__mp.td__";
export const UserDataElId = "__mp.ud__";

export type TokenData = {
    access_token?: string;
    id_token: string;
    refresh_token: string;
    expires_in: number;
    new_password_required: boolean;
    session: string;
};

export type UserData = {
    email: string;
    family_ids: string[];
    name: string;
    roles: Roles[];
    user_id: string;
    username: string;
    verified: boolean;
    exp: number;
};

export type AuthCookieBody = {
    token: TokenData;
    user: UserData;
};

type JwtData = {[key: string]: never};

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
            family_ids: [],
            name: data["name"],
            roles: data["cognito:groups"] || [],
            user_id: data["sub"],
            username: data["cognito:username"],
            verified: data["email_verified"],
            exp: data["exp"],
        };
    } catch (err) {
        console.warn(`failed to decode jwt token: ${token}`);
    }

    return undefined;
};
