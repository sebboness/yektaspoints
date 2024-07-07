"use client";

import React, { useRef } from "react";
import { PlusCircleIcon, MinusCircleIcon, CurrencyDollarIcon } from "@heroicons/react/24/solid";

import CardSingleBody from "../common/CardSingleBody";
import IssuePointsDialog, { IssuePointsDialogInterface } from "../points/IssuePointsDialog";
import { PointRequestType } from "@/lib/models/Points";

type Props = {
    childId: string;
}

const colors = {
    iconAdd: "bg-green-400 text-green-700",
    iconSubtract: "bg-red-400 text-red-700",
    iconCashout: "bg-teal-400 text-teal-700",
}

const PointActions = (props: Props) => {

    const issueDialog = useRef<IssuePointsDialogInterface>();

    const onIssuePointsClick = (type: string) => {
        console.log("issueDialog.current", issueDialog.current);
        if (issueDialog.current)
            issueDialog.current.open(type);
    };

    return (
        <CardSingleBody>
            <div className="btn-group">
                <button
                    onClick={() => onIssuePointsClick(PointRequestType.ADD)}
                    className={`btn btn-lg btn-success text-white`}>
                    <PlusCircleIcon className={`${colors.iconAdd} rounded-full h-12 w-12`} />
                    Add
                </button>
                <button
                    onClick={() => onIssuePointsClick(PointRequestType.SUBTRACT)}
                    className={`btn btn-lg btn-error text-white`}>
                    <MinusCircleIcon className={`${colors.iconSubtract} rounded-full h-12 w-12`} />
                    Subtract
                </button>
                <button
                    onClick={() => onIssuePointsClick(PointRequestType.CASHOUT)}
                    className={`btn btn-lg btn-info text-white`}>
                    <CurrencyDollarIcon className={`${colors.iconCashout} rounded-full h-12 w-12`} />
                    Cashout
                </button>
            </div>

            <IssuePointsDialog childId={props.childId} ref={issueDialog} />
        </CardSingleBody>
    );
}

export default PointActions;
