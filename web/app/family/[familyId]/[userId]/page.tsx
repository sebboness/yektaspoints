import React from "react";

import { AuthWrapper } from "@/components/AuthWrapper";
import ChildsPoints from "@/components/family/ChildsPoints";
import { cookies } from "next/headers";
import authCookie from "@/lib/auth/AuthCookie";
import { MyPointsApi } from "@/lib/api/MyPointsApi";

type Props = {
    params: {
        userId: string,
    },
};

export default async function ChildDetailPage(props: Props) {

    const token = authCookie.getTokenData(cookies());
    const api = MyPointsApi.getInstance().withToken(token?.id_token)
    const pointsResult = await api.getUserPoints(props.params.userId);
    const isSSR = typeof window === "undefined";

    if (pointsResult.status !== "SUCCESS")
        throw new Error(pointsResult.errors.join("; "));
    if (pointsResult.data === undefined || pointsResult.data === null)
        throw new Error("invalid getFamily result: " + JSON.stringify(pointsResult));

    console.log("isSSR?", isSSR);

    return (
        <AuthWrapper>
            <div className="w-screen xl:grid gap-8 xl:grid-cols-5 p-12">
                <ChildsPoints
                    childUserId={props.params.userId}
                    initialPoints={pointsResult.data.points}
                    isSSR={isSSR}
                />
            </div>
        </AuthWrapper>
    );
}
