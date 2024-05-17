import { Api, TokenGetter } from "./Api";
import { PointsList, RequestPointsRequest, RequestPointsResponse, UserPoints } from "../models/Points";
import { TokenData, UserData } from "../auth/Auth";
import { ResultT } from "./Result";

// Define base URIs for different environments
const baseUris: {[key: string]: string} = {
    "test":    "http://localhost",
    "local":   "https://mypoints-api-dev.hexonite.net",
    "dev":     "https://mypoints-api-dev.hexonite.net",
    "staging": "https://mypoints-api-staging.hexonite.net",
    "prod":    "https://api.points4us.com",
};

export class MyPointsApi extends Api {
    private static instance: MyPointsApi;

    constructor(env: string) {
        const baseUri = baseUris[env];
        super(baseUri);
        console.info(`${this.logName()}Using ${env}:${baseUri} version of api`);
        console.info(`${this.logName()}process.env`, process.env);
    }

    public static getInstance(): MyPointsApi {
        if (!MyPointsApi.instance) {
            MyPointsApi.instance = new MyPointsApi(process.env["ENV"] || "local");
        }
        return MyPointsApi.instance;
    }

    public withToken(tokenGetter: TokenGetter): MyPointsApi {
        this.tokenGetter = tokenGetter;
        return this;
    }

    ////
    // Auth
    ////

    public authenticate(username: string, password: string): Promise<ResultT<TokenData>> {
        return this.post("auth/token", {
            payload: {
                grant_type: "password",
                username,
                password,
            }
        });
    }

    public refreshToken(username: string, refresh_token: string): Promise<ResultT<TokenData>> {
        return this.post("auth/token", {
            payload: {
                grant_type: "refresh_token",
                username,
                refresh_token,
            }
        });
    }

    ////
    // User
    ////

    public getUser(): Promise<ResultT<UserData>> {
        return this.get("v1/user");
    }

    ////
    // Points
    ////
    
    public getPointsByUser(userID: string): Promise<ResultT<PointsList>> {
        return this.get(`v1/points/user/${userID}`);
    }

    public getPointSummaryByUser(userID: string): Promise<ResultT<UserPoints>> {
        return this.get(`v1/points/summary/${userID}`);
    }

    public postRequestPoints(payload: RequestPointsRequest): Promise<ResultT<RequestPointsResponse>> {
        return this.post("v1/points", {
            payload,
        });
    }
}