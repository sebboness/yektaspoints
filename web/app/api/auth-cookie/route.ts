import { NewErrorResult, NewSuccessResult } from "@/lib/api/Result";
import { NextRequest, NextResponse } from "next/server";

import { AuthCookieBody } from "@/lib/auth/Auth";
import { HttpStatus } from "@/lib/HttpStatusCodes";
import authCookie from "@/lib/auth/AuthCookie";

export async function GET(req: NextRequest) {
    const tokenData = authCookie.getTokenData(req.cookies);
    const userData = authCookie.getUserData(req.cookies);
    
    if (!tokenData) {
        return NextResponse.json(NewErrorResult("token cookie not set"), {
            status: HttpStatus.NotFound,
            statusText: "Token cookie not set",
        });
    }

    if (!userData) {
        return NextResponse.json(NewErrorResult("user cookie not set"), {
            status: HttpStatus.NotFound,
            statusText: "User cookie not set",
        });
    }

    const body: AuthCookieBody = {
        token: tokenData,
        user: userData,
    }

    return NextResponse.json(NewSuccessResult(body), { status: 200 });
}

export async function POST(req: NextRequest) {
    const body = await req.json<AuthCookieBody>();
    const domain = req.nextUrl.hostname;

    const response = NextResponse.json(NewSuccessResult(true), {
        status: HttpStatus.Created,
        statusText: "Set cookie successfully",
    });

    authCookie.setTokenData(response, domain, body.token);
    authCookie.setUserData(response, domain, body.user);

    return response;
}

export async function DELETE(req: NextRequest) {
    let response = NextResponse.json(NewSuccessResult(true), {
        status: HttpStatus.OK,
        statusText: "Auth cookie deleted successfully",
    });

    const domain = req.nextUrl.hostname;
    authCookie.deleteAll(response, domain);
    return response;
}