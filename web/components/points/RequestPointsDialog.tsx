"use client";

import * as yup from "yup";

import React, { useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "@/store/hooks";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import moment from "moment";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";

const ln = () => `[${moment().toISOString()}] RequestPointsForm: `;

const formSchema = yup.object({
    points: yup.number().integer().min(0).max(1000),
    reason: yup.string().required().min(5).max(256),
});

export const requestPointsDialogID = "request_points_dialog";

type FormData = {
    points?: number;
    reason: string;
}

const RequestPointsDialog = () => {

    const dialogRef = useRef<HTMLDialogElement>(null);

    const dispatch = useAppDispatch();
    const authState = useAppSelector((state) => state.auth);
    const api = MyPointsApi.getInstance();
    
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

    const onSubmit = async (data: FormData) => {  
        setLoading(true);      
        console.log(`${ln()}submitted data`, data);

        

        setLoading(false);
    }

    const doClose = (e: React.MouseEvent<HTMLElement>) => {
        if (dialogRef.current)
            dialogRef.current.close();

        e.preventDefault();
        return false;
    }

    return (
        <dialog id={requestPointsDialogID} className="modal" ref={dialogRef}>
            <div className="modal-box">
                <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onClick={doClose}>✕</button>
                <h3 className="font-bold text-lg">Earn some points</h3>
                <p className="py-4">Press ESC key or click the button below to close</p>
                <form method="dialog" onSubmit={handleSubmit(onSubmit)}>
                    <div className="form-control">
                        <input type="text" placeholder="Points" className="input input-bordered" { ...register("points")} />
                        <label className={`label ${errors.points ? "visible" : "invisible"}`}>
                            <a href="#" className="label-text-alt link link-hover text-red-600">Enter some points (for example "5")</a>
                        </label>
                    </div>

                    <div className="form-control">
                        <textarea className="textarea textarea-bordered" placeholder="Reason" { ...register("reason")}></textarea>
                        <label className={`label ${errors.reason ? "visible" : "invisible"}`}>
                            <a href="#" className="label-text-alt link link-hover text-red-600">Enter a reason to receive points</a>
                        </label>
                    </div>

                    <div className="modal-action">
                        <button type="submit" className={`btn btn-primary ${loading ? "btn-disabled" : ""}`}>
                            {loading
                                ? <><FontAwesomeIcon icon={faSpinner} spin /> Requesting...</>
                                : "Request"}
                        </button>

                        <button className="btn btn-secondary">Cancel</button>
                    </div>
                </form>
            </div>
        </dialog>
    );
}

export default RequestPointsDialog;
