import { DefaultStubOptions, FetchResolver } from "@/__mocks__/FetchResolver";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { NewErrorResult, NewSuccessResult } from "@/lib/api/Result";

// it("Should call dev api", async () => {
//     const api = new MyPointsApi("dev");
//     const authResult = await api.authenticate("sebboness", "sD97$$5L");
//     expect(authResult.status).toEqual("SUCCESS");
    
//     if (authResult.status === "SUCCESS") {
//         const refreshResult = await api.refreshToken("sebboness", authResult.data?.refresh_token!)
//         expect(refreshResult.status).toEqual("SUCCESS");

//         if (refreshResult.status === "SUCCESS") {
//             const getUser = await api.getUser()
//             expect(getUser.status).toEqual("SUCCESS");
//         }
//     }
// });


it("Should construct MyPointsApi", () => {
    expect(() => new MyPointsApi("")).toThrow("API: baseUri is empty, but it must be defined");
    expect(() => new MyPointsApi("abc")).toThrow("API: baseUri is empty, but it must be defined");
    expect(() => new MyPointsApi("local")).toBeDefined();
    expect(() => new MyPointsApi("dev")).toBeDefined();
    expect(() => new MyPointsApi("prod")).toBeDefined();
});