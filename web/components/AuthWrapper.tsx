import { useEffect } from "react";
import { useRouter } from "next/router";
import { getTokenFromCookie } from "@/lib/auth/Auth";
import { checkUserAuth } from "@/slices/authSlice";
import { useAppStore } from "@/store/hooks";

type Props = {
   children?: React.ReactNode;
};

export const AuthWrapper = ({ children }: Props) => {
    const { push } = useRouter();
    const store = useAppStore();

    // is user defined? if not, we know user is not logged in
    if (!store.getState().auth.user) {
        console.info("no auth data in state");
        const tokenData = getTokenFromCookie();
        if (tokenData) {
            console.info("got token data from cookie");
            // checks if token in cookie is still valid
            store.dispatch(checkUserAuth(tokenData));
        } else {
            console.info("no token data in cookie");
        }
    } else {
        console.info("auth data in state", store.getState().auth);
    }



    // const { token } = getValidAuthTokens();

    // // this query will only execute if the token is valid and the user email is not already in the redux store
    // const { error, isLoading } = useGetAuthDataQuery(
    //     { token: token || "" },
    //     {
    //         // The useGetAuthDataQuery hook will not execute the query at all if these values are falsy
    //         skip: !!userEmail || !token,
    //     }
    // );

    // if the user doesnt have a valid token, redirect to login page
    // useEffect(() => {
    //     if (!token) {
    //         push("/login");
    //         // will explain this in a moment
    //         // dispatch(logout());
    //     }
    // }, [token, push]);

    // // optional: show a loading indicator while the query is loading
    // if (isLoading) {
    //     return <div>Loading...</div>;
    // }

    return children;
};