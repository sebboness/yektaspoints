import moment from "moment";

export const formatDay_DDD_MMM_DD_hmm = (d: string): string => {
    console.log("formatDay_DDD_MMM_DD_hmm", d, moment(d, moment.ISO_8601).utc(true).format("ddd, MMM Do, h:mm a"));
    return moment(d, moment.ISO_8601).utc(true).format("ddd, MMM Do, h:mm a");
};
