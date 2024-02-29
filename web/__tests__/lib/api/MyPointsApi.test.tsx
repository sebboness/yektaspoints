import { FetchResolver } from "@/__tests__/FetchResolver";
import { MyPointsApi } from "@/lib/api/MyPointsApi";

it("Should construct MyPointsApi", () => {
    expect(() => new MyPointsApi("")).toThrow("API: baseUri is empty, but it must be defined");
    expect(() => new MyPointsApi("abc")).toThrow("API: baseUri is empty, but it must be defined");
    expect(() => new MyPointsApi("local")).toBeDefined();
    expect(() => new MyPointsApi("dev")).toBeDefined();
    expect(() => new MyPointsApi("prod")).toBeDefined();
});

describe("Should call api", () => {
    let fetchResolver!: FetchResolver;
        beforeEach(() => {
        fetchResolver = new FetchResolver();
    });

    it("should load api", async () => {
        // stub
        fetchResolver.stub( "http://localhost:8080/endpoint", "post", { id: 100 }, { created: true }, 200);
        // fetch
        fetch("http://localhost:8080/endpoint", {
                 method: "post",
                 body: JSON.stringify({ id: 100 }),
            })
            .then((response) => {
                if (response.ok) {
                    response.text().then((text) => {
                        console.log(text); // { created: true }
                        expect(text).toEqual({ created: true });
                    });
                }
        });
    });

    afterEach(() => {
    //   fetchResolver.clear();
    });
});