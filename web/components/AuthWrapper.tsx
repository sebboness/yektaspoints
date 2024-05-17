// import { redirect } from "next/navigation";
// import { getTokenFromCookie } from "@/lib/auth/Auth";
// import { checkUserAuth } from "@/slices/authSlice";
// import { useAppStore, useAppSelector } from "@/store/hooks";
// import { LocalApi } from "@/lib/api/LocalApi";

import { AuthChecker } from "./AuthChecker";
import { AuthCheckerClient } from "./AuthCheckerClient";

type Props = {
   children?: React.ReactNode;
};

export const Some = async () => {
    return (<></>);
};

export const AuthWrapper = ({ children }: Props) => {
    // // const store = useAppStore();
    // const userLoggedIn = false;
    // const authState = useAppSelector((state) => state.auth);

    // // console.log("state", store.getState());
    // console.log("authState", authState);

    // // is user defined? if not, we know user is not logged in
    // if (!authState.user) {
    //     console.info("no auth data in state");
    //     // const localApi = LocalApi.getInstance().getAuthCookie()
    //     //     .then(())
    //     const tokenData = getTokenFromCookie();
    //     if (tokenData) {
    //         console.info("got token data from cookie");
    //         // checks if token in cookie is still valid
    //         // store.dispatch(checkUserAuth(tokenData));
    //     } else {
    //         console.info("no token data in cookie");
    //         redirect("/login");
    //     }
    // } else {
    //     // console.info("auth data in state", store.getState().auth);
    // }

    // useEffect(() => {
    //     if (!userLoggedIn) {
    //         push("/login");
    //         // will explain this in a moment
    //         // dispatch(logout());
    //     }
    // }, [userLoggedIn, push]);



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

    return <>
        {children}
        <AuthChecker />
        <AuthCheckerClient />
    </>;
};