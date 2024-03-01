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

    callApi<T>(method: string, endpoint: string, options: CallOptions | undefined = undefined): Promise<ResultT<T>> {
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

            console.debug(`making request: ${method} ${url}`);

            fetch(url, reqOps)
                .then((resp) => {
                    console.log("response: ", resp);
                    resp.json()
                        .then((obj: any) => {
                            console.log("response json decoded: ", obj);
                            if (obj && obj.status) { // this means it's a formatted result object
                                resolve(obj);
                            }
                            else
                                resolve(NewErrorResultT(`Unknown error: ${JSON.stringify(obj)}`));
                        })
                    
                })
                .catch((err: any) => {
                    console.error("fetch caught error:", err);
                    if (err.errors)
                        reject(err); // Assume error object is a Result
                    else
                        reject(NewErrorResult(`Caught the following error: ${JSON.stringify(err)}`));
                });
        });
    }
}