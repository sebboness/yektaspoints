"use client";

import React, { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import moment from "moment";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

import { mapPointsToSummaries, mapSummaryToLitePoint, Point, PointRequestType, PointStatus, PointSummary } from "@/lib/models/Points";
import { getUserPoints } from "@/slices/pointsSlice";
import { useAppDispatch, useAppStore } from "@/store/hooks";

import PointsApprovalDialog, { PointsApprovalDialogInterface } from "../points/PointsApprovalDialog";
import CashoutList from "../points/CashoutList";
import PointRequestList from "../points/PointRequestList";
import PointsList from "../points/PointsList";
import SectionTitle from "../common/SectionTitle";

const ln = () => `[${moment().toISOString()}] ChildsPoints: `;

type Props = {
    childUserId: string;
    initialPoints: Point[];
};

const ChildsPoints = ({ childUserId, initialPoints }: Props) => {

    const approvalDialog = useRef<PointsApprovalDialogInterface>();
    const dispatch = useAppDispatch();
    const store = useAppStore();

    const userId = useParams()["user_id"] || "none";
    const familyId = useParams()["family_id"] || "none";

    const storePoints = store.getState().points.userPoints;
    
    const [points] = useState<Point[]>(storePoints);
    const [loading, setLoading] = useState(false);

    const settledPoints = points.filter(x => x.status === PointStatus.SETTLED && x.request.type !== PointRequestType.CASHOUT);
    const requestedPoints = points.filter(x => x.status === PointStatus.WAITING);
    const cashouts = points.filter(x => x.status === PointStatus.SETTLED && x.request.type === PointRequestType.CASHOUT);

    const handleOnRequestClick = (p: PointSummary) => {
        const point = mapSummaryToLitePoint(p);
        point.user_id = userId;
        console.log("point", point);
        console.log("approvalDialog.current", approvalDialog.current);
        approvalDialog.current?.open(point);
    };

    useEffect(() => {
        setLoading(true);
        console.log(`${ln()}dispatching getUserPoints`);

        dispatch(getUserPoints(userId));
        setLoading(false);

        // const fetchedPoints = store.getState().points.userPoints;
        // if (fetchedPoints) {
        //     console.log(`${ln()}setting user points`, fetchedPoints);
        //     setPoints(fetchedPoints);
        // };
    }, []);

    return (
        <>
            {/* Left */}
            <div className="container mx-auto col-span-3">
                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500 mb-8">
                    <div className="card-body">
                        <SectionTitle>
                            [Name]&apos;s points&nbsp;
                            {loading
                                ? <FontAwesomeIcon icon={faSpinner} spin />
                                : <></>}
                        </SectionTitle>

                        <PointsList points={mapPointsToSummaries(settledPoints)} />
                    </div>
                </div>
            </div>

            {/* Right */}
            <div className="container mx-auto col-span-2">
                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 mb-16 border border-zinc-500">
                    <div className="card-body">
                        <SectionTitle>
                            [Name]&apos;s requests&nbsp;
                            {loading
                                ? <FontAwesomeIcon icon={faSpinner} spin />
                                : <></>}
                        </SectionTitle>

                        <PointRequestList
                            onClick={(p) => handleOnRequestClick(p)}
                            points={mapPointsToSummaries(requestedPoints)} />
                    </div>
                </div>

                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                    <div className="card-body">
                        <SectionTitle>
                            [Name]&apos;s cashout history&nbsp;
                            {loading
                                ? <FontAwesomeIcon icon={faSpinner} spin />
                                : <></>}
                        </SectionTitle>

                        <CashoutList points={mapPointsToSummaries(cashouts)} />
                    </div>
                </div>
            </div>

            <PointsApprovalDialog ref={approvalDialog} />
        </>
    );
};

export default ChildsPoints;
