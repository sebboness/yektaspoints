import React from "react";
import { AuthChecker } from "./AuthChecker";
import { AuthCheckerClient } from "./AuthCheckerClient";

type Props = {
   children?: React.ReactNode;
};

export const AuthWrapper = ({ children }: Props) => {
    return <>
        <AuthChecker />
        <AuthCheckerClient />
        {children}
    </>;
};