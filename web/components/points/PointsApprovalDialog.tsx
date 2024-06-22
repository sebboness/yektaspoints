"use client";

import moment from "moment";
import React, { useEffect, useImperativeHandle, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { yupResolver } from "@hookform/resolvers/yup";

// import pointsSlice, { PointsSlice } from "@/slices/pointsSlice";
// import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { Point, PointDecisionApprove, PointDecisionDeny } from "@/lib/models/Points";
import { getTokenRetriever } from "@/store/store";
import { FamilyMember } from "@/lib/models/Family";

const ln = () => `[${moment().toISOString()}] PointsApprovalDialog: `;

const formSchema = yup.object({
    decision: yup.string().nullable().oneOf([PointDecisionApprove, PointDecisionDeny]),
    point_id: yup.string().required().min(0).max(1000),
    parent_notes: yup.string().required().min(5).max(256),
});

export const approvePointsRequestDialogID = "approve_points_request_dialog";

export interface PointsApprovalDialogInterface {
    open(point: Point, child: FamilyMember): void;
}

type FormData = {
    point_id: string;
    decision?: string | null;
    parent_notes: string | null;
};

const PointsApprovalDialog = React.forwardRef((props, ref) => {

    const [mounted, setMounted] = useState(false);

    const dialogRef = useRef<HTMLDialogElement>(null);

    useImperativeHandle(ref, () => ({
        close: () => close(),
        open: (point: Point, child: FamilyMember) => open(point, child),
    }));

    const api = MyPointsApi.getInstance();
    
    const [loading, setLoading] = useState(false);
    const [decision, setDecision] = useState("");
    const [point, setPoint] = useState<Point|undefined>(undefined);
    const [child, setChild] = useState<FamilyMember|undefined>(undefined);

    // Setup form validation variables and methods
    const { 
        register,
        handleSubmit,
        reset,
        // watch,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(formSchema)
    });

    const onSubmit = async (data: FormData) => {  
        setLoading(true);      
        console.log(`${ln()}submitted data`, data);

        if (!point)
            return;

        // const result = await api
        //     .withToken(getTokenRetriever())
        //     .approveRequestPoints({
        //         decision: data.decision || "",
        //         parent_notes: data.parent_notes,
        //         point_id: point.id,
        //         user_id: point.user_id,
        //     });

        // if (result.status === "SUCCESS") {
        //     console.log(`${ln()}approve/deny point request`, result);
        //     // dispatch(PointsSlice.actions.addPointToRequesting(result.data.point_summary));
        //     close();
        // } else {
        //     console.log(`${ln()}error approve/deny point request`, result);
        // }

        // setLoading(false);
    };

    const close = () => {
        reset();

        setPoint(undefined);
        setChild(undefined);

        if (dialogRef.current)
            dialogRef.current.close();
    };

    const doClose = (e: React.MouseEvent<HTMLElement>) => {
        close();

        e.preventDefault();
        return false;
    };

    const open = (point: Point, child: FamilyMember) => {
        console.log("opened approval dialog", point);
        if (dialogRef.current) {
            setPoint(point);
            setChild(child);
            dialogRef.current.showModal();
        }
    };
    
    // Ensure component is mounted
    useEffect(() => setMounted(true), []);
    if (!mounted) {
        return null;
    }

    return (
        <dialog id={approvePointsRequestDialogID} className="modal" ref={dialogRef}>
            <div className="modal-box bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2" onClick={doClose}>✕</button>
                

                <div className="divide-y divide-blue-200">
                    <div>
                        <h3 className="font-bold text-lg">{child ? child.name : ""} requested pointes</h3>
                        <p className="py-4 text-center text-2xl">
                            {point
                                ? `${point.points} point${point.points == 1 ? "" : "s"}`
                                : ""}
                        </p>
                        <p className="text-lg">
                            <strong>Reason: </strong>
                            {point
                                ? (point.request.reason || "No reason given")
                                : ""}
                        </p>
                    </div>

                    <div>
                        <form method="dialog" onSubmit={handleSubmit(onSubmit)}>
                            <input type="hidden" { ...register("point_id")} value={point ? point.id : ""}/>
                            <input type="hidden" { ...register("decision")} value={decision}/>
                            <div className="form-control">
                                <textarea className="textarea textarea-bordered" placeholder="Optional: Notes for [NAME]" { ...register("parent_notes")}></textarea>
                                <label className={`label ${errors.parent_notes ? "visible" : "invisible"}`}>
                                    <a href="#" className="label-text-alt link link-hover text-red-600"></a>
                                </label>
                            </div>

                            <div className="modal-action">
                                {loading
                                    ? <><FontAwesomeIcon icon={faSpinner} spin /> Saving&hellip;</>
                                    : <></>}
                                <>
                                    <button
                                        className={`btn btn-primary ${loading ? "btn-disabled" : ""}`}
                                        onClick={() => {
                                            setDecision(PointDecisionApprove);
                                            
                                        }}
                                    >Approve</button>
                                    <button
                                        className={`btn btn-primary ${loading ? "btn-disabled" : ""}`}
                                        onClick={() => setDecision(PointDecisionDeny)}
                                    >Deny</button>
                                    <a
                                        className={`btn btn-secondary ${loading ? "btn-disabled" : ""}`}
                                        onClick={doClose}>
                                        Cancel
                                    </a>
                                </>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </dialog>
    );
});

PointsApprovalDialog.displayName = "PointsApprovalDialog";

export default PointsApprovalDialog;
