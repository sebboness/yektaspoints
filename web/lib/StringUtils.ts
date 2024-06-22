/**
 * Generates and returns a string of random alpha-numeric characters of given length
 * @param length Positive integer greater than one
 * @returns Random string of alpha-numeric characters
 */
export const randomAlphaNumericString = (length: number): string => {
    if (length <= 1)
        throw "length must be a positive integer greater than 1";

    let outString: string = "";
    let inOptions: string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

    for (let i = 0; i < length; i++) {
        outString += inOptions.charAt(Math.floor(Math.random() * inOptions.length));
    }

    return outString;
}