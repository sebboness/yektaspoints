import React from "react";
import moment from "moment";
import { cookies } from "next/headers";

import ParentsApp from "@/components/app/ParentsApp";
import authCookie from "@/lib/auth/AuthCookie";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { ThrowIfNotSuccess } from "@/lib/api/Result";
import { MapType } from "@/lib/models/Common";
import { Family } from "@/lib/models/Family";
import { AuthWrapper } from "@/components/AuthWrapper";

const ln = () => `[${moment().toISOString()}] ParentsAppLanding: `;

const ParentsAppLanding = async () => {

    const tokenData = authCookie.getTokenData(cookies());
    const api = MyPointsApi.getInstance().withToken(tokenData?.id_token);

    console.log(`${ln()}token? ${tokenData ? (tokenData.id_token.substring(tokenData.id_token.length - 20)) : "NONE"}`);
    console.log(`${ln()}B`);

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
        <AuthWrapper>
            <ParentsApp user={user} families={families} />
        </AuthWrapper>
    );
};

export default ParentsAppLanding;
