import { HttpStatus } from "@/lib/HttpStatusCodes";

test.skip('Workaround', () => {});

type Method = "get" | "options" | "post" | "put" | "patch" | "delete";

export const DefaultStubOptions: StubOptions = {
    reqHeaders: {
        "Content-Type": "json/application",
    },
    credentials: "include",
}

export type StubOptions = {
    reqHeaders: HeadersInit | undefined
    credentials: RequestCredentials | undefined
}

/**
 * Stub API request, response in test cases.
 * - should be initialized and destroyed within the context of a specific case.
 * - highly customizable
 *
 * <pre>
 *  describe("Fetch API", () => {
 *    let fetchResolver!: FetchResolver;
 *      beforeEach(() => {
 *      fetchResolver = new FetchResolver();
 *    });
 *
 *    it("should load api", async () => {
 *      // stub
 *      fetchResolver.stub( "http://localhost:8080/endpoint", "post", { id: 100 }, { created: true }, 200);
 *      // fetch
 *      fetch("http://localhost:8080/endpoint",
 *        { method: "post", body: JSON.stringify({ id: 100 })}
 *      ).then((response) => {
 *        if (response.ok) {
 *          response.text().then((text) => {
 *            console.log(text); // { created: true }
 *            expect(text).toBeEqual({ created: true });
 *          });
 *        }
 *      });
 *    });
 *
 *    afterEach(() => {
 *      fetchResolver.clear();
 *    });
 *  });
 * </pre>
 *
 * Even though jest executes tests in parallel jest instance,
 * We can't go wrong if stubs are cleaned after its use
 */
export class FetchResolver {
    private mocks: Map<string, Response> = new Map();

    constructor() {
        this.init();
    }

    public stub(uri: string, method: Method, payload: any, response: any, status: HttpStatus, opts: StubOptions | undefined = undefined) {

        const finalRequest: { input: RequestInfo | URL; init?: RequestInit } = {
            input: uri,
            init: {
                method: method,
                body: JSON.stringify(payload),
                headers: opts?.reqHeaders,
                credentials: opts?.credentials,
            }
        };

        console.log(
            `mocking fetch :::\nrequest ${this.prettyPrint(
            finalRequest
            )} with \nresponse ${this.prettyPrint(response)} ans status ${status}`
        );
        
        this.mocks.set(
            JSON.stringify(finalRequest),
            new Response(JSON.stringify(response), { status: status })
        );
    }

    private prettyPrint(json: any) {
        return JSON.stringify(json, null, 2);
    }

    public clear() {
        this.mocks.clear();
    }

    private init() {
        jest
            .spyOn(global, "fetch")
            .mockImplementation((input: RequestInfo | URL, init?: RequestInit) => {
                const request = {input, init };

                return new Promise((resolve, reject) => {
                    const key = JSON.stringify(request);
                    const response = this.mocks.get(key);

                    if (response) {
                        resolve(response);
                    }
                    else {
                        // rejecting here will hurt component initialization
                        console.error(
                            `mock not implemented :::\nrequest ${this.prettyPrint(request)}`
                        );
                        
                        // return empty response
                        resolve(new Response("{}"));
                    }
                });
        });
    }
}