"use client";

import * as yup from "yup";

import React, { useRef, useState } from "react";
import pointsSlice, { PointsSlice } from "@/slices/pointsSlice";
import { useAppDispatch, useAppSelector } from "@/store/hooks";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { Point } from "@/lib/models/Points";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { getTokenRetriever } from "@/store/store";
import moment from "moment";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";

const ln = () => `[${moment().toISOString()}] PointsApprovalDialog: `;

const formSchema = yup.object({
    points: yup.number().integer().min(0).max(1000),
    notes: yup.string().required().min(5).max(256),
});

export const requestPointsDialogID = "request_points_dialog";

type Props = {
    point: Point
}

type FormData = {
    pointId: string;
    decision: string;
    notes?: string;
}

const PointsApprovalDialog = (props: Props) => {

    const dialogRef = useRef<HTMLDialogElement>(null);

    const dispatch = useAppDispatch();
    const authState = useAppSelector((state) => state.auth);
    const api = MyPointsApi.getInstance();
    
    const [loading, setLoading] = useState(false);
    const [decision, setDecision] = useState("");

    // Setup form validation variables and methods
    const { 
        register,
        handleSubmit,
        reset,
        watch,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(formSchema)
    });

    const onSubmit = async (data: FormData) => {  
        setLoading(true);      
        console.log(`${ln()}submitted data`, data);

        // const result = await api
        //     .withToken(getTokenRetriever())
        //     .postRequestPoints({
        //         points: data.points || 0,
        //         notes: data.notes,
        //     })

        // if (result.data) {
        //     console.log(`${ln()}approve/deny point request`, result.data);
        //     dispatch(PointsSlice.actions.addPointToRequesting(result.data.point_summary));
        //     close();
        // } else {
        //     console.log(`${ln()}error approve/deny point request`, result);
        // }

        setLoading(false);
    }

    const close = () => {
        reset();

        if (dialogRef.current)
            dialogRef.current.close();
    }

    const doClose = (e: React.MouseEvent<HTMLElement>) => {
        close();

        e.preventDefault();
        return false;
    }

    return (
        <dialog id={requestPointsDialogID} className="modal" ref={dialogRef}>
            <div className="modal-box bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onClick={doClose}>✕</button>
                

                <div className="divide-y divide-blue-200">
                    <div>
                        <h3 className="font-bold text-lg">[NAME] requested pointes</h3>
                        <p className="py-4 text-center text-2xl">10 points</p>
                        <p className="text-lg">[REASON...]</p>
                    </div>

                    <div>
                        {/* <form method="dialog" onSubmit={handleSubmit(onSubmit)}>
                            <input type="hidden" { ...register("pointId")} value={props.point.id}/>
                            <input type="hidden" { ...register("decision")} value={decision}/>
                            <div className="form-control">
                                <textarea className="textarea textarea-bordered" placeholder="Optional: Notes for [NAME]" { ...register("notes")}></textarea>
                                <label className={`label ${errors.notes ? "visible" : "invisible"}`}>
                                    <a href="#" className="label-text-alt link link-hover text-red-600"></a>
                                </label>
                            </div>

                            <div className="modal-action">
                                {loading
                                    ? <>
                                        <button className={`btn btn-primary ${loading ? "btn-disabled" : ""}`} onClick={setDecision("APPROVE")}>Approve</button>
                                        <button className={`btn btn-primary ${loading ? "btn-disabled" : ""}`} onClick={setDecision("DENY ")}>Deny</button>
                                        <a className="btn btn-secondary" onClick={doClose}>Cancel</a>
                                    </>
                                    : <>Loading...</>}
                            </div>
                        </form> */}
                    </div>
                </div>
            </div>
        </dialog>
    );
}

export default PointsApprovalDialog;
