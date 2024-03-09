import { NewErrorResult, NewSuccessResult } from "@/lib/api/Result";
import { NextRequest, NextResponse } from "next/server";

import { TokenData } from "@/lib/auth/Auth";
import authCookie from "@/lib/auth/AuthCookie";

export async function GET(req: NextRequest) {
    const tokenData = await authCookie.getTokenData(req.cookies);
    if (!tokenData) {
        return NextResponse.json(NewErrorResult("not set"), {
            status: 404,
            statusText: "Auth cookie not set",
        });
    }

    return NextResponse.json(NewSuccessResult(tokenData), { status: 200 });
}

export async function POST(req: NextRequest) {
    const body = await req.json<TokenData>();
    const domain = req.nextUrl.hostname;

    const response = NextResponse.json(NewSuccessResult(true), {
        status: 201,
        statusText: "Set cookie successfully",
    });

    await authCookie.setTokenData(response.cookies, domain, body);

    return response;
}

export async function DELETE(req: NextRequest) {
    let response = NextResponse.json(NewSuccessResult(true), {
        status: 200,
        statusText: "Auth cookie deleted successfully",
    });

    const domain = req.nextUrl.hostname;
    authCookie.deleteAll(response, domain);
    return response;
}