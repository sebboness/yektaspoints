import React from "react";
import { cookies } from "next/headers";
import moment from "moment";

import { AuthWrapper } from "@/components/AuthWrapper";
import KidsList from "@/components/family/KidsList";
import authCookie from "@/lib/auth/AuthCookie";
import { MyPointsApi } from "@/lib/api/MyPointsApi";

const ln = () => `[${moment().toISOString()}] Family: `;

type Props = {
    params: {
        familyId: string,
    },
};

export default async function Family(props: Props) {

    const token = authCookie.getTokenData(cookies());
    console.log(`${ln()}token? ${token ? (token.id_token.substring(token.id_token.length - 20)) : "NONE"}`);
    console.log(`${ln()}B`);

    const api = MyPointsApi.getInstance().withToken(token?.id_token)
    const familyResult = await api.getFamily(props.params.familyId);

    if (familyResult.status !== "SUCCESS")
        throw new Error(familyResult.errors.join("; "));
    if (!familyResult.data)
        throw new Error("invalid getFamily result: " + JSON.stringify(familyResult));

    return (
        <AuthWrapper>
            <section>
                <div className="w-screen gap-8 grid grid-cols-1 p-12">
                    <div className="container mx-auto">
                        <KidsList familyId={props.params.familyId} initialFamily={familyResult.data.family} />
                    </div>
                </div>
            </section>
        </AuthWrapper>
    );
}
