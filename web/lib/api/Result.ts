export type Result = {
    status: "FAILURE" | "SUCCESS";
    errors: string[];
    message: string | null;
};

export type ResultT<T> = Result & {
    data: T | null;
};
