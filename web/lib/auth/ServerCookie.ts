"use server";

import { ResponseCookie } from "next/dist/compiled/@edge-runtime/cookies";
import { cookies } from "next/headers";

export const setCookieServerSide = (cookieOptions: ResponseCookie) => {
    cookies().set(cookieOptions);
}