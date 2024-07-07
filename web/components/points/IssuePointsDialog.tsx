"use client";

import moment from "moment";
import React, { useEffect, useImperativeHandle, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { PointsSlice } from "@/slices/pointsSlice";
import { useAppDispatch } from "@/store/hooks";
import { getTokenRetriever } from "@/store/store";
import { PointRequestType } from "@/lib/models/Points";

const ln = () => `[${moment().toISOString()}] IssuePointsDialog: `;

const formSchema = yup.object({
    parent_notes: yup.string().max(500),
    points: yup.number().integer().min(0).max(1000),
    reason: yup.string().required().min(5).max(256),
    to_user_id: yup.string().required().max(36),
    type: yup.string().required().oneOf([PointRequestType.ADD, PointRequestType.CASHOUT, PointRequestType.SUBTRACT])
});

export const issuePointsDialogID = "issue_points_dialog";

type FormData = {
    parent_notes?: string;
    points?: number;
    reason: string;
    to_user_id: string;
    type: string;
};

type Props = {
    childId: string;
}

export interface IssuePointsDialogInterface {
    open(type: string): void;
}

const IssuePointsDialog = React.forwardRef((props: Props, ref) => {

    const [mounted, setMounted] = useState(false);
    const [type, setType] = useState("");

    const dialogRef = useRef<HTMLDialogElement>(null);

    useImperativeHandle(ref, () => ({
        close: () => close(),
        open: (type: string) => open(type),
    }));

    const dispatch = useAppDispatch();
    const api = MyPointsApi.getInstance();
    
    const [loading, setLoading] = useState(false);

    // Setup form validation variables and methods
    const { 
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(formSchema)
    });

    const onSubmit = async (data: FormData) => {  
        setLoading(true);      
        console.log(`${ln()}submitted data`, data);

        const result = await api
            .withToken(getTokenRetriever())
            .issuePoints({
                parent_notes: data.parent_notes || null,
                points: data.points || 0,
                reason: data.reason,
                to_user_id: props.childId,
                type: type,
            });

        if (result.data) {
            console.log(`${ln()}issue points response`, result.data);
            dispatch(PointsSlice.actions.addPointToRequesting(result.data.point_summary));
            close();
        } else {
            console.log(`${ln()}error requesting points`, result);
        }

        setLoading(false);
    };

    const close = () => {
        setType("");
        reset();

        if (dialogRef.current)
            dialogRef.current.close();
    };

    const doClose = (e: React.MouseEvent<HTMLElement>) => {
        close();

        e.preventDefault();
        return false;
    };

    const open = (type: string) => {
        console.log("open", type);
        setType(type);
        if (dialogRef.current)
            dialogRef.current.showModal();
    };
    
    // Ensure component is mounted
    useEffect(() => setMounted(true), []);
    if (!mounted) {
        return null;
    }

    return (
        <dialog id={issuePointsDialogID} className="modal" ref={dialogRef}>
            <div className="modal-box bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onClick={doClose}>✕</button>
                
                <h3 className="font-bold text-lg">Earn some points</h3>

                <p className="py-4">Press ESC key or click the button below to close</p>

                <form method="dialog" onSubmit={handleSubmit(onSubmit)}>
                    <div className="form-control">
                        <input type="number" placeholder="Points" className="input input-bordered" { ...register("points")} />
                        <label className={`label ${errors.points ? "visible" : "invisible"}`}>
                            <a href="#" className="label-text-alt link link-hover text-red-600">Enter positive points (for example 5)</a>
                        </label>
                    </div>

                    <div className="form-control">
                        <textarea className="textarea textarea-bordered" placeholder="Reason" { ...register("reason")}></textarea>
                        <label className={`label ${errors.reason ? "visible" : "invisible"}`}>
                            <a href="#" className="label-text-alt link link-hover text-red-600">Enter a reason to receive points</a>
                        </label>
                    </div>

                    <div className="form-control">
                        <textarea className="textarea textarea-bordered" placeholder={`Optional notes`} { ...register("parent_notes")}></textarea>
                        <label className={`label ${errors.parent_notes ? "visible" : "invisible"}`}>
                            <a href="#" className="label-text-alt link link-hover text-red-600"></a>
                        </label>
                    </div>

                    <div className="modal-action">
                        <button type="submit" className={`btn btn-primary ${loading ? "btn-disabled" : ""}`}>
                            {loading
                                ? <><FontAwesomeIcon icon={faSpinner} spin /> Requesting...</>
                                : "Request"}
                        </button>

                        <a className="btn btn-secondary" onClick={doClose}>Cancel</a>
                    </div>
                </form>
            </div>
        </dialog>
    );
});

IssuePointsDialog.displayName = "IssuePointsDialog";

export default IssuePointsDialog;
