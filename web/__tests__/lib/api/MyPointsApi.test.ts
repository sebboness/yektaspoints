import { DefaultStubOptions, FetchResolver } from "@/__tests__/FetchResolver";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { NewErrorResult, NewSuccessResult, ResultT } from "@/lib/api/Result";

it("Should construct MyPointsApi", () => {
    expect(() => new MyPointsApi("")).toThrow("API: baseUri is empty, but it must be defined");
    expect(() => new MyPointsApi("abc")).toThrow("API: baseUri is empty, but it must be defined");
    expect(() => new MyPointsApi("local")).toBeDefined();
    expect(() => new MyPointsApi("dev")).toBeDefined();
    expect(() => new MyPointsApi("prod")).toBeDefined();
});

describe("Should call api", () => {
    let fetchResolver: FetchResolver;
    beforeEach(() => {
        fetchResolver = new FetchResolver();
    });

    type RT = {
        age: number
    }

    it("should call successfully", async () => {
        // Arrange
        const respJson = NewSuccessResult({ age: 123 })
        fetchResolver.stub( "http://localhost/endpoint", "post", undefined, respJson, 200, DefaultStubOptions);

        const api = new MyPointsApi("test");
        const result = await api.callApi<RT>("post", "endpoint");
        expect(result.status).toEqual("SUCCESS");
        expect(result.data).toBeDefined();
        expect(result.data?.age).toEqual(123);
    });

    it("should fail properly with a server error", async () => {
        // Arrange
        const respJson = NewErrorResult("internal server error", "fail");
        fetchResolver.stub( "http://localhost/endpoint", "get", undefined, respJson, 500, DefaultStubOptions);

        const api = new MyPointsApi("test");
        const result = await api.callApi<RT>("get", "endpoint");
        expect(result.message).toEqual("fail");
        expect(result.status).toEqual("FAILURE");
        expect(result.errors).toContain("internal server error");
        expect(result.data).toBeNull();
    });

    afterEach(() => {
        if (fetchResolver)
            fetchResolver.clear();
    });
});