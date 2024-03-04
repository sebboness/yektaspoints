import { NewErrorResult, NewSuccessResult } from "@/lib/api/Result";
import { TokenData } from "@/lib/auth/Auth";
import authCookie from "@/lib/auth/AuthCookie";
import { NextRequest, NextResponse } from "next/server";


export async function GET(req: NextRequest) {
    const tokenData = await authCookie.get();
    if (!tokenData) {
        return NextResponse.json(NewErrorResult("not set"), {
            status: 404,
            statusText: "Auth cookie not set",
        });
    }

    // const cookieDecompressed = await decompress(cookie.value);
    // const tokenData = JSON.parse(cookieDecompressed) as TokenData;

    return NextResponse.json(NewSuccessResult(tokenData), { status: 200 });
}

export async function POST(req: NextRequest) {
    const body = await req.json<TokenData>();
    const domain = req.nextUrl.hostname;

    const response = NextResponse.json(NewSuccessResult(true), {
        status: 201,
        statusText: "Set cookie successfully",
    });

    await authCookie.set(response, domain, body);

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