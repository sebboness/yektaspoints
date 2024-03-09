import AuthMiddleware from "@/middleware/AuthMiddleware";

export const xRedirectToHeader = "x-middleware-request-redirect";

export const middleware = [
    AuthMiddleware,
];
