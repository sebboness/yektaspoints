"use client";

import React from "react";

import { useAppSelector } from "@/store/hooks";

const Greeter = () => {

    // const { push } = useRouter();
    const user = useAppSelector ((state) => state.auth.user);
    // const [displayTxt, setDisplayTxt] = useState("No user")

    // useEffect(() => {

    //     if (!user) {
    //         setDisplayTxt("No user");
    //     } else if (!user.roles) {
    //         setDisplayTxt("User has no roles");
    //     } else if (user.roles.findIndex((x) => x === Roles.Parent)) {
    //         push("/mykids");
    //     } else if (user.roles.findIndex((x) => x === Roles.Child)) {
    //         push("/mypoints");
    //     } else if (user.roles.findIndex((x) => x === Roles.Admin)) {
    //         push("/admin");
    //     }

    // }, [user]);

    return (<>
        {/* <p>{displayTxt}</p> */}
        <p>{JSON.stringify(user)}</p>
    </>);
};

export default Greeter;
