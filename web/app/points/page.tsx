import React from "react";

import { AuthWrapper } from "@/components/AuthWrapper";
import UserSummary from "@/components/points/UserSummary";

export default async function Home() {

    return (
        <AuthWrapper>
            {/* Top navbar */}
            {/* TK TK */}
        
            <section>
                <UserSummary />
            </section>
        </AuthWrapper>
    );
}
