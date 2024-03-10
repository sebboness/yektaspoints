import { NextRequest, NextResponse } from "next/server";

import authCookie from "@/lib/auth/AuthCookie";
import moment from "moment";
import { xRedirectToHeader } from "./config";

const ln = () => `[${moment().toISOString()}] LogoutMiddleware: `;

export default async function HomeMiddleware(req: NextRequest): Promise<NextResponse> {
    
    const path = req.nextUrl.pathname;
    
    if (path != "/logout") {
        return NextResponse.next();
    } else {
        console.log(`${ln()}path match on ${path}: ${req.nextUrl.href}`);

        const res = NextResponse.next();

        // delete all auth-related cookies
        authCookie.deleteAll(res, req.nextUrl.hostname);

        // redirect to home page
        res.headers.set(xRedirectToHeader, "/"); 

        return res;
    }
}