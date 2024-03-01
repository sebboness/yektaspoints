import { Api } from "./Api";
import { NewErrorResult, NewErrorResultT, ResultT, SUCCESS } from "./Result";

// Define base URIs for different environments
const baseUris: {[key: string]: string} = {
    "test":    "http://localhost",
    "local":   "https://mypoints-api-dev.hexonite.net",
    "dev":     "https://mypoints-api-dev.hexonite.net",
    "staging": "https://mypoints-api-staging.hexonite.net",
    "prod":    "https://mypoints-api.hexonite.net",
};

type QueryParams = {[key: string]: string[]} | undefined;

type CallOptions = {
    withAuth: boolean;
    queryParams: QueryParams;
    payload: object | undefined;
};

export class MyPointsApi extends Api {
    private static instance: MyPointsApi;

    constructor(env: string) {
        if (env !== "prod") {
            console.info(`MyPointsApi: Using ${env} version of api`);
        }

        super(baseUris[env]);
    }

    public static getInstance(): MyPointsApi {
        if (!MyPointsApi.instance) {
            MyPointsApi.instance = new MyPointsApi(process.env["ENV"] || "local");
        }
        return MyPointsApi.instance;
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
            const headers: HeadersInit = {
                "Content-Type": "json/application",
            };

            if (options?.withAuth) {
                headers["Authentication"] = "";
            }

            // Build call url
            const url = this.getCallUrl(this.baseUri, endpoint, options?.queryParams);

            // Initialize fetch request
            const reqOps: RequestInit = {
                method,
                headers,
                credentials: "include",
            };

            // Add payload
            if (options?.payload)
                reqOps.body = JSON.stringify(options.payload);

            console.debug(`MyPointsAPi making request: ${method} ${url}`);

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
                    if (err && err.errors)
                        resolve(err); // Assume error object is a Result
                    else if (typeof err === "string" || err instanceof Array)
                        resolve(NewErrorResultT(err, "Caught an error"));
                    else
                        resolve(NewErrorResultT(`Caught the following error: ${JSON.stringify(err)}`));
                });
        });
    }

    public post<T>(endpoint: string, options: CallOptions | undefined = undefined): Promise<ResultT<T>> {
        return this.callApi<T>("post", endpoint, options);
    }
}