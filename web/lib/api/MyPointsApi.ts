import { Api, TokenGetter } from "./Api";
import { Result, ResultT } from "./Result";
import { TokenData, UserData } from "../auth/Auth";
import { ApprovePointsRequest, PointsList, RequestPointsRequest, RequestPointsResponse, UserPoints } from "../models/Points";
import { FamilyResponse } from "../models/Family";

// Define base URIs for different environments
const baseUris: {[key: string]: string} = {
    "test":    "http://localhost",
    "local":   "http://localhost:10010",
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
        console.info(`${this.logName()}process.env.ENV`, process.env.ENV);
        console.info(`${this.logName()}process.env.NEXT_PUBLIC_ENV`, process.env.NEXT_PUBLIC_ENV);
        console.info(`${this.logName()}process.env.NEXT_PUBLIC_ABC`, process.env.NEXT_PUBLIC_ABC);
    }

    public static getInstance(): MyPointsApi {
        if (!MyPointsApi.instance) {
            MyPointsApi.instance = new MyPointsApi(process.env.NEXT_PUBLIC_ENV || "local");
        }
        return MyPointsApi.instance;
    }

    public withToken(tokenOrGetter?: TokenGetter | string): MyPointsApi {
        if (typeof tokenOrGetter === "string")
            this.tokenGetter = {
                getToken: () => tokenOrGetter,
                getTokenType: () => "Bearer",
            };
        else
            this.tokenGetter = tokenOrGetter;
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
    // Family
    ////

    public getFamily(familyId: string): Promise<ResultT<FamilyResponse>> {
        return this.get(`v1/family/${familyId}`);
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

    public postRequestPoints(payload: RequestPointsRequest): Promise<ResultT<RequestPointsResponse>> {
        return this.post("v1/points ", {
            payload,
        });
    }

    public approveRequestPoints(payload: ApprovePointsRequest): Promise<Result> {
        return this.put(`v1/points/${payload.point_id}/approve`, {
            payload,
        });
    }

    ////
    // User points
    ////
    
    public getUserPoints(userID: string): Promise<ResultT<PointsList>> {
        return this.get(`v1/user/${userID}/points`);
    }

    public getUserPointsSummary(userID: string): Promise<ResultT<UserPoints>> {
        return this.get(`v1/user/${userID}/points-summary`);
    }
}