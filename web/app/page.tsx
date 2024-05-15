import React from "react";

import { AuthWrapper } from "@/components/AuthWrapper";
import Greeter from "@/components/Greeter";

export default async function Home() {

    return (
        <AuthWrapper>
            <Greeter />
        </AuthWrapper>
    );
}
