import { TokenData } from "../auth/Auth";
import { Api } from "./Api";
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

    public getAuthCookie(): Promise<ResultT<TokenData>> {
        return this.get("api/auth-cookie");
    }

    public setAuthCookie(tokenData: TokenData): Promise<ResultT<boolean>> {
        return this.post("api/auth-cookie", {
            payload: tokenData,
        });
    }
}