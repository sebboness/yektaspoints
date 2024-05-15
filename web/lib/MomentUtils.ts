import moment from "moment";

export const formatDay_DDD_MMM_DD_hmm = (d: string): string => {
    return moment(d, moment.ISO_8601).utc(true).local().format("ddd, MMM Do, h:mm a");
};
