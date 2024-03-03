"use client";

import React, { FormEvent, useState } from "react"
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCompass } from "@fortawesome/free-solid-svg-icons";

const formSchema = yup.object({
    username: yup.string().required().max(30),
    password: yup.string().required().max(256),
});

type FormBody = {
    username?: string;
    password?: string;
}

const LoginForm = () => {
    // Loading when credentials are submitted
    const [loading, setLoading] = useState(false);

    // Setup form validation variables and methods
    const { 
        register,
        handleSubmit,
        watch,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(formSchema)
    });

    const onSubmit = (data: FormBody) => {  
        setLoading(true);      
        console.log("on submit", data);
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