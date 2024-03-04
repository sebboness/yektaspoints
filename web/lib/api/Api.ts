import { NewErrorResultT, ResultT } from "./Result";

export type QueryParams = {[key: string]: string[]} | undefined;

export type CallOptions = {
    payload?: any;
    queryParams?: QueryParams;
    token?: string;
};

/**
 * Interface that defines an object from which an auth token is retrieved.
 * This is used when calling the API to determine if the Authorized header should be set
 */
export interface TokenGetter {
    getToken: () => string;
    getTokenType: () => string;
}

export class Api {
    baseUri: string;
    tokenGetter?: TokenGetter;

    constructor(baseUri: string) {
        this.baseUri = baseUri;
    }

    private getCallUrl(baseUri: string, endpoint: string, queryParams: QueryParams): string {
        let queryString = "";
        if (queryParams) {
            const parts: Array<string> = [];
            for (let [k, v] of Object.entries(queryParams)) {
                const key = encodeURIComponent(k);
                const val = encodeURIComponent(v.join(","));
                parts.push(`${key}=${val}`)
            }
            queryString = "?" + parts.join("&");
        }

        return `${baseUri}/${endpoint}${queryString}`;
    }

    public callApi<T>(method: string, endpoint: string, options: CallOptions | undefined = undefined): Promise<ResultT<T>> {
        return new Promise((resolve, reject) => {
            const opts = options || {};
            const headers: HeadersInit = {
                "Content-Type": "json/application",
                // "Origin": "http://localhost:3000",
            };

            let isAuthedReq = false;

            // Attach auth token to headers if it is set
            if (this.tokenGetter) {
                headers["Authentication"] = this.tokenGetter.getTokenType() + " " + this.tokenGetter.getToken();
                isAuthedReq = true;
            }

            // Build call url
            const url = this.getCallUrl(this.baseUri, endpoint, opts.queryParams);

            // Initialize fetch request
            const reqOps: RequestInit = {
                method,
                headers,
                credentials: "include",
            };

            // Add payload
            if (opts.payload)
                reqOps.body = JSON.stringify(opts.payload);

            console.info(`MyPointsAPI preparing${isAuthedReq ? " authorized" : ""} request: ${method} ${url}`);

            fetch(url, reqOps)
                .then((resp) => {
                    console.log("MyPointsAPi response: ", resp);
                    resp.json()
                        .then((obj: any) => {
                            console.log("response json decoded: ", obj);
                            if (obj && obj.status) { // this means it's a formatted result object
                                resolve(obj);
                            }
                            else
                                resolve(NewErrorResultT(`Unexpected response: ${JSON.stringify(obj)}`));
                        })
                    
                })
                .catch((err: any) => {
                    console.error("MyPointsAPi caught an error:", err);
                    console.error("MyPointsAPi typeof error:", typeof err);
                    if (err && err.errors)
                        resolve(err); // Assume error object is a Result
                    else if (typeof err === "string" || err instanceof Array)
                        resolve(NewErrorResultT(err, "Caught an error"));
                    else
                        resolve(NewErrorResultT(`Caught the following error: ${err}`));
                });
        });
    }

    public delete<T>(endpoint: string, opts: CallOptions | undefined = undefined): Promise<ResultT<T>> {
        return this.callApi<T>("delete", endpoint, opts);
    }

    public get<T>(endpoint: string, opts: CallOptions | undefined = undefined): Promise<ResultT<T>> {
        return this.callApi<T>("get", endpoint, opts);
    }

    public patch<T>(endpoint: string, opts: CallOptions | undefined = undefined): Promise<ResultT<T>> {
        return this.callApi<T>("patch", endpoint, opts);
    }

    public post<T>(endpoint: string, opts: CallOptions | undefined = undefined): Promise<ResultT<T>> {
        return this.callApi<T>("post", endpoint, opts);
    }

    public put<T>(endpoint: string, opts: CallOptions | undefined = undefined): Promise<ResultT<T>> {
        return this.callApi<T>("put", endpoint, opts);
    }
}