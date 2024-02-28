import { Api } from "./Api";
import { ResultT } from "./Result";

// Define base URIs for different environments
const baseUris: {[key: string]: string} = {
    "local":   "https://mypoints-api-dev.hexonite.net",
    "dev":     "https://mypoints-api-dev.hexonite.net",
    "staging": "https://mypoints-api-staging.hexonite.net",
    "prod":    "https://mypoints-api.hexonite.net",
};

type CallOptions = {
    withAuth: boolean;
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

    callApi<T>(method: string, endpoint: string, options: CallOptions | undefined): Promise<ResultT<T>> {
        return new Promise((resolve, reject) => {
            const headers: HeadersInit = {
                "Content-Type": "json/application",
            };

            if (options?.withAuth) {
                headers["Authentication"] = "";
            }

            const url = `${this.baseUri}/${endpoint}`;

            const reqOps: RequestInit = {
                method,
                headers,
                credentials: "include",
            };

            fetch(url, reqOps)
                .then((resp) => {
                    resp.json();
                })
                .then((obj: ResultT<T>) => {
                    
                })
                .catch((err) => {

                });
        });
    }
}