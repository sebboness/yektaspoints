import { NextRequest, NextResponse } from "next/server";
import { middleware as activatedMiddleware, xRedirectToHeader } from "@/middleware/config";

import moment from "moment";

const ln = () => `[${moment().toISOString()}] middleware: `;

export async function middleware(req: NextRequest) {
    // Initialize a NextResponse object
    const nextResponse = NextResponse.next();

    // Map through activated middleware functions
    const middlewareFunctions = activatedMiddleware.map(fn => fn(req));

    // Array to store middleware headers
    const middlewareHeader = [];

    // Loop through middleware functions
    for (const middleware of middlewareFunctions) {

        // Execute middleware function and await the result
        const result = await middleware;

        // Check if the result is not okay and return it
        if (!result.ok) {
            return result;
        }
        // Push middleware headers to the array
        middlewareHeader.push(result.headers);
    }

    //First we are going to define a redirectTo variable
    let redirectTo = null;

    // Check each header in middlewareHeader
    middlewareHeader.some((header) => {
        // Look for the 'x-middleware-request-redirect' header
        const redirect = header.get(xRedirectToHeader);

        console.log(`${ln()}redirect set? ${redirect}`)
        
        if (redirect) {
            // If a redirect is found, store the value and break the loop
            redirectTo = redirect;
            return true; // Break the loop
        }
        // Continue to the next header in case the redirect header is not found
        return false; // Continue the loop
    });

    // If a redirection is required based on the middleware headers
    if (redirectTo) {
        // Perform the redirection
        return NextResponse.redirect(new URL(redirectTo, req.url), {
            status: 307, // Use the appropriate HTTP status code for the redirect
        });
    }

    // If no redirection is needed, proceed to the next middleware or route handler
    return nextResponse;
}