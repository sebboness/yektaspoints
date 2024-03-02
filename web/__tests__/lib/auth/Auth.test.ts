import { ParseToken } from "@/lib/auth/Auth";

describe("Decode JWT token", () => {
    it("Should decode token", () => {
        const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjMiLCJjb2duaXRvOmdyb3VwcyI6WyJhZG1pbiIsInBhcmVudCJdLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiY29nbml0bzp1c2VybmFtZSI6ImpvaG4iLCJuYW1lIjoiSm9obiIsImV4cCI6MTcwOTQxNTczOCwiaWF0IjoxNzA5NDEyMTM4LCJqdGkiOiJhNTgxNmVkNi01NGRmLTRiMmMtYmIxMi00MzllNTE2MTZjY2YiLCJlbWFpbCI6ImpvaG5AaW5mby5jbyJ9.Aja5K-0U7SeJXuYwYkuhUv-iXKys8Rx3m0N2k2gKR0I";
        const user = ParseToken(token);
        expect(user).toBeDefined();
        expect(user?.email).toEqual("john@info.co");
        expect(user?.name).toEqual("John");
        expect(user?.groups).toEqual(["admin","parent"]);
        expect(user?.user_id).toEqual("123");
        expect(user?.username).toEqual("john");
        expect(user?.verified).toEqual(true);
        expect(user?.exp).toEqual(1709415738);
    });

    it("Should fail to decode token", () => {
        const token = "blah";
        const user = ParseToken(token);
        expect(user).toBeUndefined();
    });
});