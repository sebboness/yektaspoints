"use client";

import React from "react";
import { useAppSelector  } from "@/store/hooks";
import { Roles } from "@/lib/auth/Roles";
import { useRouter } from "next/navigation";

const Greeter = () => {

    const { push } = useRouter();
    const user = useAppSelector ((state) => state.auth.user);

    if (!user) {
        return (<p>No user</p>);
    } else if (!user.roles) {
        return (<p>No roles</p>);
    } else if (user.roles.findIndex((x) => x === Roles.Parent)) {
        push("/mykids");
        return;
    } else if (user.roles.findIndex((x) => x === Roles.Child)) {
        push("/mypoints");
        return;
    } else if (user.roles.findIndex((x) => x === Roles.Admin)) {
        push("/admin");
        return;
    }

    return (<>
        <p>Hello! I am {user.name}</p>
        <p>My family IDs: {user.family_ids}</p>
        <p>My roles: {user.roles}</p>
    </>);
}

export default Greeter