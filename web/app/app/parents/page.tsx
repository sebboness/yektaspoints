import React from "react";
import moment from "moment";
import { headers } from "next/headers";

import ParentsApp from "@/components/app/ParentsApp";
import authCookie, { TokenHeaderName } from "@/lib/auth/AuthCookie";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { ThrowIfNotSuccess } from "@/lib/api/Result";
import { MapType } from "@/lib/models/Common";
import { Family } from "@/lib/models/Family";

const ln = () => `[${moment().toISOString()}] ParentsAppLanding: `;

const ParentsAppLanding = async () => {

    const tokenData = authCookie.getTokenDataFromHeader(headers());
    const token = tokenData ? tokenData.id_token : "none";
    console.log(`${ln()}token? ${token ? (token.substring(token.length - 20)) : "NONE"}`);

    const api = MyPointsApi.getInstance().withToken(token);
    const userResult = await api.getUser();
    ThrowIfNotSuccess(userResult);

    const user = userResult.data!;
    const families: MapType<Family> = {};

    for (const familyId of user.family_ids) {
        const familyResult = await api.getFamily(familyId);
        if (familyResult.data) {
            families[familyId] = familyResult.data.family;
        }
    }

    return (
        <ParentsApp user={user} families={families} />
    );
};

export default ParentsAppLanding;
