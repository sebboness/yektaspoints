export const SUCCESS = "SUCCESS";
export const FAILURE = "FAILURE";

export type Result = {
    status: "FAILURE" | "SUCCESS";
    errors: string[];
    message?: string | null;
    statusCode?: number;
};

export type ResultT<T> = Result & {
    data: T | null;
};

export function NewErrorResult(err: string | Array<string>, msg: string | null = null, statusCode: number | undefined = undefined): Result {
    return NewErrorResultT<any>(err, msg, statusCode);
};

export function NewErrorResultT<T>(err: string | Array<string>, msg: string | null = null, statusCode: number | undefined = undefined): ResultT<T> {
    return {
        status: FAILURE,
        errors: typeof err === "string" ? [err] : err,
        message: msg || null,
        data: null,
        statusCode: statusCode,
    }
};

export function NewSuccessResult<T> (data: T, msg: string | null = null, statusCode: number | undefined = undefined): ResultT<T> {
    return {
        status: SUCCESS,
        errors: [],
        message: msg || null,
        data: data,
        statusCode: statusCode,
    }
};


/**
 * Returns an error result containing information of the given error object
 * @param err The error object which could be a result object or any other kind of error
 * @returns The result object with error information
 */
export const ErrorAsResult = <T>(err: ResultT<T> | any): ResultT<T> => {
    if (err.status)
        return err;
    else
        return NewErrorResultT(JSON.stringify(err));
}