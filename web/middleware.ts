import { NextRequest, NextResponse } from "next/server";
import { middleware as activatedMiddleware, xRedirectToHeader } from "@/middleware/config";

import { HttpStatus } from "./lib/HttpStatusCodes";
import { ResponseCookie } from "next/dist/compiled/@edge-runtime/cookies";
import moment from "moment";

const ln = () => `[${moment().toISOString()}] middleware: `;

export async function middleware(req: NextRequest) {
    // Initialize a NextResponse object
    let res = NextResponse.next();

    // Map through activated middleware functions
    const middlewareFunctions = activatedMiddleware.map(fn => fn(req));

    // Array to store middleware headers
    const middlewareHeader = [];
    const middlewareCookies: ResponseCookie[] = [];

    // Loop through middleware functions
    for (const middleware of middlewareFunctions) {

        // Execute middleware function and await the result
        const result = await middleware;

        // Check if the result is not okay and return it
        if (!result.ok) {
            return result;
        }

        // Push middleware headers and cookies to the header array
        middlewareHeader.push(result.headers);

        // Push new cookies into the cookie array
        result.cookies.getAll().forEach((c) => {
            console.log(`${ln()}checking cookie ${c.name}`);
            if (!middlewareCookies.some((mc) => mc.name === c.name)) {
                console.log(`${ln()}pushing cookie ${c.name}`);
                middlewareCookies.push(c);
            }
        });        

        const redirect = result.headers.get(xRedirectToHeader);
        if (redirect) {
            // redirect is present, so let's quit the middleware loop
            console.log(`${ln()}redirect found in middleware. quitting loop`);
            break;
        }
    }

    // First we are going to define a redirectTo variable
    let redirectTo = null;

    // Check each header in middlewareHeader
    middlewareHeader.some((header) => {
        // Look for the 'x-middleware-request-redirect' header
        const redirect = header.get(xRedirectToHeader);

        
        if (redirect) {
            // If a redirect is found, store the value and break the loop
            console.log(`${ln()}redirect set to ${redirect}`)
            redirectTo = redirect;
            return true; // Break the loop
        }
        // Continue to the next header in case the redirect header is not found
        return false; // Continue the loop
    });

    // If a redirection is required based on the middleware headers
    if (redirectTo) {
        res = NextResponse.redirect(new URL(redirectTo, req.url), {
            status: HttpStatus.TemporaryRedirect, // Use the appropriate HTTP status code for the redirect
        });
    }

    // Set response cookies
    middlewareCookies.forEach((c) => {
        console.log(`${ln()}setting cookie ${c.name}`);
        res.cookies.set(c);
    });

    // If no redirection is needed, proceed to the next middleware or route handler
    return res;
}