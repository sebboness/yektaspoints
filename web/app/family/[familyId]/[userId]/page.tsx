import React from "react";

import { AuthWrapper } from "@/components/AuthWrapper";
import ChildsPoints from "@/components/family/ChildsPoints";

type Props = {
    params: {
        userId: string,
    },
};

export default async function ChildDetailPage(props: Props) {

    return (
        <AuthWrapper>
            <div className="w-screen xl:grid gap-8 xl:grid-cols-5 p-12">
                <ChildsPoints childUserId={props.params.userId} />
            </div>
        </AuthWrapper>
    );
}
