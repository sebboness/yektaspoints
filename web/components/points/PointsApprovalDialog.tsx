"use client";

import moment from "moment";
import React, { useEffect, useImperativeHandle, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { yupResolver } from "@hookform/resolvers/yup";

// import pointsSlice, { PointsSlice } from "@/slices/pointsSlice";
import { MyPointsApi } from "@/lib/api/MyPointsApi";
import { Point, PointDecisionApprove, PointDecisionDeny } from "@/lib/models/Points";
import { FamilyMember } from "@/lib/models/Family";
import { PointsSlice } from "@/slices/pointsSlice";
import { useAppDispatch } from "@/store/hooks";

const ln = () => `[${moment().toISOString()}] PointsApprovalDialog: `;

const formSchema = yup.object({
    decision: yup.string().nullable().oneOf([PointDecisionApprove, PointDecisionDeny]),
    point_id: yup.string().required().min(0).max(36),
    parent_notes: yup.string().max(500),
});

export const approvePointsRequestDialogID = "approve_points_request_dialog";

export interface PointsApprovalDialogInterface {
    open(point: Point, child: FamilyMember): void;
}

type FormData = {
    point_id: string;
    decision?: string | null;
    parent_notes?: string;
};

type State = {
    point: Point;
    child: FamilyMember;
}

const PointsApprovalDialog = React.forwardRef((props, ref) => {

    const [mounted, setMounted] = useState(false);

    const dispatch = useAppDispatch();
    const dialogRef = useRef<HTMLDialogElement>(null);
    const formRef = useRef<HTMLFormElement>(null);

    useImperativeHandle(ref, () => ({
        close: () => close(),
        open: (point: Point, child: FamilyMember) => open(point, child),
    }));

    // const api = MyPointsApi.getInstance();
    
    const [loading, setLoading] = useState(false);
    const [decision, setDecision] = useState("");
    const [data, setData] = useState<State|undefined>(undefined);

    // Setup form validation variables and methods
    const { 
        register,
        handleSubmit,
        reset,
        // watch,
        setValue,
        formState: { errors },
    } = useForm({
        resolver: yupResolver(formSchema)
    });

    const onSubmit = async (formData: FormData) => {
        if (!data) {
            // TODO trigger warning flash
            console.warn(`${ln()}no data in form?`);
            close();
            return;
        }

        setLoading(true);      
        console.log(`${ln()}submitted data`, formData);

        // dispatch(PointsSlice.actions.onApproveRequesting(data.point));

        const api = MyPointsApi.getInstance();
        const result = await api.approveRequestPoints({
                decision: formData.decision || "",
                parent_notes: formData.parent_notes || null,
                point_id: formData.point_id,
                user_id: data.point.user_id || "",
            });

        if (result.status === "SUCCESS" && result.data) {
            console.log(`${ln()}approve/deny point request`, result);
            dispatch(PointsSlice.actions.onApproveRequesting(result.data.point));
            close();
        } else {
            console.log(`${ln()}error approve/deny point request`, result);
        }

        setLoading(false);
    };

    const close = () => {
        reset();

        setData(undefined);
        setLoading(false);

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
            setData({ point, child });
            setValue("point_id", point.id);
            dialogRef.current.showModal();
        }
    };

    const childName = data ? data.child.name : "";

    useEffect(() => {
        setValue("decision", decision);

        // console.log(decision, formRef.current);
        // if (decision && formRef.current) {
        //     // formRef.current.submit();
        //     console.log("hello", errors);
        // }
    }, [decision]);
    
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
                        <p className="py-4 text-2xl">
                            {childName} is requesting&nbsp;
                            <span className="font-bold">
                                {data
                                    ? `${data.point.points} point${data.point.points == 1 ? "" : "s"}`
                                    : ""}
                            </span>.
                        </p>
                        <p className="text-lg">
                            <span className="font-bold">Reason:</span>
                            <br />
                            {data
                                ? (data.point.request.reason || "No reason given")
                                : ""}
                        </p>
                    </div>

                    <div>
                        <p className="text-lg">
                            <span className="font-bold">Respond with optional notes:</span>
                        </p>
                        <form method="dialog" onSubmit={handleSubmit(onSubmit)} ref={formRef}>
                            <div className="form-control">
                                <textarea className="textarea textarea-bordered" placeholder={`Optional: Notes for ${childName}`} { ...register("parent_notes")}></textarea>
                                <label className={`label ${errors.parent_notes ? "visible" : "invisible"}`}>
                                    <a href="#" className="label-text-alt link link-hover text-red-600"></a>
                                </label>
                            </div>

                            <div className="modal-action place-content-between">
                                {loading
                                    ? <div>
                                        <button className="btn btn-disabled"><FontAwesomeIcon icon={faSpinner} spin /> Saving&hellip;</button>
                                    </div>
                                    : <div className="grid grid-cols-2 gap-2">
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
                                    </div>}
                                <div>
                                    <a
                                        className={`btn btn-secondary ${loading ? "btn-disabled" : ""}`}
                                        onClick={doClose}>
                                        Cancel
                                    </a>
                                </div>
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
