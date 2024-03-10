import AuthMiddleware from "@/middleware/AuthMiddleware";
import HomeMiddleware from "@/middleware/HomeMiddleware";
import LogoutMiddleware from "@/middleware/LogoutMiddleware";

export const xRedirectToHeader = "x-middleware-request-redirect";

export const middleware = [
    LogoutMiddleware,
    AuthMiddleware,
    HomeMiddleware,
];
