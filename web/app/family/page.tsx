import { AuthWrapper } from "@/components/AuthWrapper";
import CardSingleBody from "@/components/common/CardSingleBody";
import KidsList from "@/components/family/KidsList";
import React from "react";

export default async function Family() {

    return (
        <AuthWrapper>        
            <section>
                <div className="w-screen gap-8 grid grid-cols-1 p-12">
                    <div className="container mx-auto">
                        <CardSingleBody>
                            <KidsList />
                        </CardSingleBody>
                    </div>                    
                </div>
            </section>
        </AuthWrapper>
    );
}
