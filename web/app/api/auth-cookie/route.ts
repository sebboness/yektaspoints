import { NewErrorResult, NewSuccessResult } from "@/lib/api/Result";
import { TokenData } from "@/lib/auth/Auth";
import { NextRequest, NextResponse } from "next/server";

const cookieName = "app_auth";

export async function GET(req: NextRequest) {
    const cookie = req.cookies.get(cookieName);
    console.info(`${cookieName} set? ${cookie !== undefined}`);
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

    const response = NextResponse.json(NewSuccessResult(true), {
        status: 201,
        statusText: "Set cookie successfully",
    });

    console.info(`Setting ${cookieName} cookie with value ${JSON.stringify(body)}`);
    response.cookies.set({
        name: cookieName,
        value: JSON.stringify(body),
        maxAge: 60*60*24*30, // 30 days
        httpOnly: true,
        sameSite: "strict",
    });

    return response;
}

export async function DELETE(req: NextRequest) {
    const response = NextResponse.json(NewSuccessResult(true), {
        status: 200,
        statusText: "Auth cookie deleted successfully",
    });

    console.info(`Deleting ${cookieName} cookie`);
    response.cookies.set({
        name: cookieName,
        value: "",
        maxAge: -100,
        httpOnly: true,
        sameSite: "strict",
    });

    return response;
}