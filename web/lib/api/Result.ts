export const SUCCESS = "SUCCESS";
export const FAILURE = "FAILURE";

export type Result = {
    status: "FAILURE" | "SUCCESS";
    errors: string[];
    message: string | null;
};

export type ResultT<T> = Result & {
    data: T | null;
};

export function NewErrorResult(err: string | Array<string>, msg: string | null = null): Result {
    return NewErrorResultT<any>(err, msg);
};

export function NewErrorResultT<T>(err: string | Array<string>, msg: string | null = null): ResultT<T> {
    return {
        status: FAILURE,
        errors: typeof err === "string" ? [err] : err,
        message: msg || null,
        data: null,
    }
};

export function NewSuccessResult<T> (data: T, msg: string | null = null): ResultT<T> {
    return {
        status: SUCCESS,
        errors: [],
        message: msg || null,
        data: data,
    }
};