import { Api } from "./Api";
import { AuthCookieBody } from "../auth/Auth";
import { ResultT } from "./Result";

export class LocalApi extends Api {
    private static instance: LocalApi;

    constructor() {
        super("");
    }

    public static getInstance(): LocalApi {
        if (!LocalApi.instance) {
            LocalApi.instance = new LocalApi();
        }
        return LocalApi.instance;
    }

    public deleteAuthCookie(): Promise<ResultT<boolean>> {
        return this.delete("api/auth-cookie");
    }

    public getAuthCookie(): Promise<ResultT<AuthCookieBody>> {
        return this.get("api/auth-cookie");
    }

    public setAuthCookie(body: AuthCookieBody): Promise<ResultT<boolean>> {
        return this.post("api/auth-cookie", {
            payload: body,
        });
    }
}