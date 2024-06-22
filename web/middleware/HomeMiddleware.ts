import { NextRequest, NextResponse } from "next/server";

import { Roles } from "@/lib/auth/Roles";
import authCookie from "@/lib/auth/AuthCookie";
import moment from "moment";
import { xRedirectToHeader } from "./config";

const ln = () => `[${moment().toISOString()}] HomeMiddleware: `;

export default async function HomeMiddleware(req: NextRequest): Promise<NextResponse> {
    
    const path = req.nextUrl.pathname;
    
    if (path != "/") {
        return NextResponse.next();
    } else {
        console.log(`${ln()}path match on ${path}: ${req.nextUrl.href}`);

        const response = NextResponse.next();
        const user = authCookie.getUserData(req.cookies);

        if (!user) {
            console.log(`${ln()}no user cookie available, so not logged in`);
            return response;
        }

        console.log(`${ln()}user ${user.username}, roles: ${user.roles}`);

        if (user.roles.some((r) => r === Roles.Parent)) {
            // User is parent, so redirect to family landing page
            response.headers.set(xRedirectToHeader, "/app/parents");
        } else if (user.roles.some((r) => r === Roles.Child)) {
            // User is child, so redirect to points landing page
            response.headers.set(xRedirectToHeader, "/app/mypoints");
        } else if (user.roles.some((r) => r === Roles.Admin)) {
            // User is child, so redirect to points landing page
            response.headers.set(xRedirectToHeader, "/admin");
        }
        
        console.log(`${ln()}${response.headers.has(xRedirectToHeader) ? ("redirect to " + response.headers.get(xRedirectToHeader)) : "no redirects"}`);    

        return response;
    }
}