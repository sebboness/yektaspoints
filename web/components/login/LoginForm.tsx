"use client";

import React, { useState } from "react"
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCompass } from "@fortawesome/free-solid-svg-icons";
import { TokenData } from "@/lib/auth/Auth";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { SUCCESS } from "@/lib/api/Result";
import { useAppDispatch } from "@/store/hooks";
import authSlice, { AuthSlice, setAuthCookie } from "@/slices/authSlice";

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

    const dispatch = useAppDispatch();
    const api = MyPointsApi.getInstance();

    // api.callApi("GET", "health", {})
    //     .then((res) => {
    //         console.info("health res:", res);
    //     })
    //     .catch((err) => {
    //         console.error("health err:", err);
    //     });
    
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
        console.log("on submit", data);

        const result = await api.authenticate(data.username!, data.password!);
        if (result.status === SUCCESS) {
            console.log("api logged in", result.data);            
            dispatch(AuthSlice.actions.setAuthToken(result.data!));
            dispatch(setAuthCookie(result.data!));
            return
        }

        console.log("api error", result);
        setLoading(false);
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
                        ? <><FontAwesomeIcon icon={faCompass} spin /> Logging in...</>
                        : "Login"}
                </button>
            </div>
        </form>
    )
}

LoginForm.propTypes = {}

export default LoginForm