import React from "react";
import { cookies } from "next/headers";

import { AuthWrapper } from "@/components/AuthWrapper";
import authCookie from "@/lib/auth/AuthCookie";
import FamilyList from "@/components/family/FamilyList";

export default async function Family() {

    const userData = authCookie.getUserData(cookies());
    const familyIds = userData
        ? userData.family_ids
        : [];

    return (
        <AuthWrapper>        
            <section>
                <div className="w-screen gap-8 grid grid-cols-1 p-12">
                    <div className="container mx-auto">
                        <FamilyList initialFamilyIds={familyIds} />
                    </div>                    
                </div>
            </section>
        </AuthWrapper>
    );
}
