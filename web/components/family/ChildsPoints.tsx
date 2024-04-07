"use client";

import {
    CircleStackIcon,
    CurrencyDollarIcon,
    HandThumbUpIcon,
    LightBulbIcon,
    SparklesIcon,
} from "@heroicons/react/24/solid";
import { PointRequestType, PointStatus, mapPointsToSummaries } from "@/lib/models/Points";
import React, { useEffect, useState } from "react";
import { getUserPointSummary, getUserPoints } from "@/slices/pointsSlice";
import { useAppDispatch, useAppSelector } from "@/store/hooks";

import CashoutList from "../points/CashoutList";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Image from "next/image";
import PointRequestList from "../points/PointRequestList";
import PointsList from "../points/PointsList";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import moment from "moment";

const ln = () => `[${moment().toISOString()}] ChildsPoints: `;

type Props = {
    childUserId: string,
}

const ChildsPoints = (props: Props) => {
    
    const [loading, setLoading] = useState(true);
    const dispatch = useAppDispatch();

    const points = useAppSelector((state) => state.points.userPoints);
    const childUserId = props.childUserId;

    const settledPoints = points.filter(x => x.status === PointStatus.SETTLED && x.request.type !== PointRequestType.CASHOUT);
    const requestedPoints = points.filter(x => x.status === PointStatus.WAITING);
    const cashouts = points.filter(x => x.status === PointStatus.SETTLED && x.request.type === PointRequestType.CASHOUT);

    useEffect(() => {
        if (childUserId) {
            console.log(`${ln()}dispatching getUserPoints`);
            dispatch(getUserPoints(childUserId));
            setLoading(false);
        }
    }, [childUserId]);

    console.log(`${ln()}info`, props.childUserId, loading, childUserId);

    return (
        <>
            {/* Left */}
            <div className="container mx-auto col-span-3">
                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500 mb-8">
                    <div className="card-body">
                        <p className="text-2xl font-bold">Child's points</p>

                        <PointsList points={mapPointsToSummaries(settledPoints)} />
                    </div>
                </div>
            </div>

            {/* Right */}
            <div className="container mx-auto col-span-2">
                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 mb-16 border border-zinc-500">
                    <div className="card-body">
                        <p className="text-2xl font-bold">Child's requests</p>

                        <PointRequestList points={mapPointsToSummaries(requestedPoints)} />
                    </div>
                </div>

                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                    <div className="card-body">
                        <p className="text-2xl font-bold">Child's cashout history</p>

                        <CashoutList points={mapPointsToSummaries(cashouts)} />
                    </div>
                </div>
            </div>
        </>
    );
}

export default ChildsPoints;
