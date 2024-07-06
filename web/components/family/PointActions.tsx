"use client";

import React from "react";
import { PlusCircleIcon, MinusCircleIcon, CurrencyDollarIcon } from "@heroicons/react/24/solid";

import CardSingleBody from "../common/CardSingleBody";

type Props = {
    childId: string;
}

const colors = {
    iconAdd: "bg-green-400 text-green-700",
    iconSubtract: "bg-red-400 text-red-700",
    iconCashout: "bg-teal-400 text-teal-700",
}

const PointActions = (props: Props) => {
    return (
        <CardSingleBody>
            <div className="btn-group">
                <button
                    className={`btn btn-lg btn-success text-white`}>
                    <PlusCircleIcon className={`${colors.iconAdd} rounded-full h-12 w-12`} />
                    Add
                </button>
                <button
                    className={`btn btn-lg btn-error text-white`}>
                    <MinusCircleIcon className={`${colors.iconSubtract} rounded-full h-12 w-12`} />
                    Subtract
                </button>
                <button
                    className={`btn btn-lg btn-info text-white`}>
                    <CurrencyDollarIcon className={`${colors.iconCashout} rounded-full h-12 w-12`} />
                    Cashout
                </button>
            </div>
        </CardSingleBody>
    );
}

export default PointActions;
