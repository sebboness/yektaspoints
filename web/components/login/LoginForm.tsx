"use client";

import * as yup from "yup";

import { AuthSlice, getSimpleTokenRetriever, setAuthCookie } from "@/slices/authSlice";
import React, { useState } from "react"
import { useRouter, useSearchParams } from "next/navigation";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import moment from "moment";
import { useAppDispatch } from "@/store/hooks";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";

const ln = () => `[${moment().toISOString()}] LoginForm: `;

const formSchema = yup.object({
    username: yup.string().required().max(30),
    password: yup.string().required().max(256),
});

export type LoginFormData = {
    username?: string;
    password?: string;
}

type Props = {}

const LoginForm = (props: Props) => {

    // router/nav stuff
    const router = useRouter();
    const searchParams = useSearchParams();
    
    const dispatch = useAppDispatch();
    const api = MyPointsApi.getInstance();
    
    // Loading when credentials are submitted
    const [loading, setLoading] = useState(false);
    const [preparing, setPreparing] = useState(false);

    // Setup form validation variables and methods
    const { 
        register,
        handleSubmit,
        watch,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(formSchema)
    });

    const onSubmit = async (data: LoginFormData) => {  
        setLoading(true);      
        console.log(`${ln()}logging in...`);

        const authResult = await api.authenticate(data.username!, data.password!);
        if (authResult.data) {
            console.log(`${ln()}api logged in. get user...`);
            setPreparing(true);
            
            const userResult = await api
                .withToken(getSimpleTokenRetriever(authResult.data.id_token))
                .getUser();

            if (userResult.data) {
                console.log(`${ln()}got user a-ok`, userResult.data);  

                dispatch(AuthSlice.actions.setAuthToken(authResult.data));
                dispatch(setAuthCookie({
                    token: authResult.data,
                    user: userResult.data,
                }));

                // redirect to where the user came from (defaults to home page)
                const returnUrl = searchParams.has("return_url")
                    ? (searchParams.get("return_url") || "/")
                    : "/";

                router.push(returnUrl);
                
                return;
            }
        }

        console.log(`${ln()}login error`, authResult);
        setLoading(false);
        setPreparing(false);
    }

    return (
        <form className="grid grid-cols-1 gap-y-4" onSubmit={handleSubmit(onSubmit)}>
            <div className="form-control">
                <input placeholder="Username" className="input input-bordered" { ...register("username")} />
            </div>
            <div className="form-control">
                <input type="password" placeholder="Password" className="input input-bordered" { ...register("password")} />
                <label className={`label ${errors.username || errors.password ? "visible" : "invisible"}`}>
                    <a href="#" className="label-text-alt link link-hover text-red-600">Username and Password are required fields</a>
                </label>
            </div>
            <div className="form-control">
                <label className="label">
                    <a href="#" className="label-text-alt link link-hover">Forgot password?</a>
                </label>
                <button type="submit" className={`btn btn-primary ${loading ? "btn-disabled" : ""}`}>
                    {loading
                        ? <><FontAwesomeIcon icon={faSpinner} spin /> Logging in...</>
                        : "Login"}
                </button>
            </div>
        </form>
    )
}

LoginForm.propTypes = {}

export default LoginForm