import React from "react";
import moment from "moment";
import { AuthWrapper } from "@/components/AuthWrapper";
import SectionTitle from "@/components/common/SectionTitle";
import Link from "next/link";
import authCookie from "@/lib/auth/AuthCookie";
import { cookies } from "next/headers";
import { MyPointsApi } from "@/lib/api/MyPointsApi";

const ln = () => `[${moment().toISOString()}] TestPageWithId: `;

type Props = {
    params: {
        id: string
    }
};

const TestPageWithId = async (props: Props) => {
    console.log(`${ln()}Hello!`);

    const token = authCookie.getTokenData(cookies());
    console.log(`${ln()}token? ${token ? (token.id_token.substring(token.id_token.length - 20)) : "NONE"}`);
    console.log(`${ln()}B`);

    const api = MyPointsApi.getInstance().withToken(token?.id_token)
    const result = await api.getUser();
    if (result.status !== "SUCCESS")
        throw new Error(result.errors.join("; "));
    if (!result.data)
        throw new Error("invalid getUser result: " + JSON.stringify(result));

    return (
        <>
            <AuthWrapper>
                <section>
                    <div className="w-screen gap-8 grid grid-cols-1 p-12">
                        <div className="container mx-auto">
                            <SectionTitle>
                                <Link href={"/family/test"}>&laquo; Back</Link>
                                <br />
                                Hello TestPage with ID {props.params.id}
                                <br />
                                User's email: {result.data.email}
                            </SectionTitle>
                        </div>
                    </div>
                </section>
            </AuthWrapper>
        </>
    );
}

export default TestPageWithId;
