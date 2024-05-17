import { Suspense } from "react";
import Image from "next/image";

import LoginForm from "@/components/login/LoginForm";
// import { redirect, useRouter } from "next/navigation";
// import { useAppSelector } from "@/store/hooks";

export default function LoginPage() {    
    
    // const authState = useAppSelector((state) => state.auth);
    // if (authState.user && authState.user.user_id) {
    //     console.info("Already logged in");
    //     redirect("/");
    // }

    return (
        <>
            {/* Top navbar */}
            {/* TK TK */}
        
            <div className="hero min-h-screen">
                <div className="card card-compact bg-base-100 shadow-xl p-4 bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                    <div className="card-body">

                        <figure className="mb-8">
                            <Image
                                src="/img/logo-512.svg"
                                width={384}
                                height={164}
                                priority={true}
                                alt="Picture of the author"
                            />
                        </figure>

                        <h1 className="text-3xl">
                            {/* {authState.user ? "Logged in" : "Not logged in"} */}
                            {/* Hello there! :) */}
                        </h1>

                        <Suspense>
                            <LoginForm />
                        </Suspense>
                    </div>
                </div>
            </div>
        </>
    );
}
