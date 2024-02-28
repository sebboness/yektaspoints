export class Api {
    baseUri: string;

    constructor(baseUri: string) {
        if (!baseUri) {
            throw "API: baseUri is empty, but it must be defined";
        }        

        this.baseUri = baseUri;
    }
}