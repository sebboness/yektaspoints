import { NewErrorResult, NewSuccessResult } from "@/lib/api/Result";
import { TokenCookieName, TokenData } from "@/lib/auth/Auth";
import { NextRequest, NextResponse } from "next/server";


const env = process.env["ENV"];

export async function GET(req: NextRequest) {
    const cookie = req.cookies.get(TokenCookieName);
    console.info(`${TokenCookieName} set? ${cookie !== undefined}`);
    if (!cookie) {
        return NextResponse.json(NewErrorResult("not set"), {
            status: 404,
            statusText: "Auth cookie not set",
        });
    }

    const tokenData = JSON.parse(cookie.value) as TokenData;
    return NextResponse.json(NewSuccessResult(tokenData), { status: 200 });
}

export async function POST(req: NextRequest) {
    const body = await req.json<TokenData>();
    const domain = req.nextUrl.hostname;

    const response = NextResponse.json(NewSuccessResult(true), {
        status: 201,
        statusText: "Set cookie successfully",
    });

    console.info(`Setting ${TokenCookieName} cookie on ${env}:${domain} with value ${JSON.stringify(body)}`);
    response.cookies.set({
        name: TokenCookieName,
        value: JSON.stringify(body),
        maxAge: 60*60*24*30, // 30 days
        httpOnly: true,
        sameSite: "strict",
        secure: env == "local" ? false : true,
        domain: domain === "localhost" ? domain : "hexonite.net",
    });

    return response;
}

export async function DELETE(req: NextRequest) {
    const response = NextResponse.json(NewSuccessResult(true), {
        status: 200,
        statusText: "Auth cookie deleted successfully",
    });

    console.info(`Deleting ${TokenCookieName} cookie`);
    response.cookies.set({
        name: TokenCookieName,
        value: "",
        maxAge: -100,
        httpOnly: true,
        sameSite: "strict",
    });

    return response;
}