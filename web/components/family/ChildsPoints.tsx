"use client";

import React, { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import moment from "moment";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

import { mapPointsToSummaries, mapSummaryToLitePoint, PointRequestType, PointStatus, PointSummary } from "@/lib/models/Points";
import { getUserPoints } from "@/slices/pointsSlice";
import { useAppDispatch, useAppSelector } from "@/store/hooks";

import PointsApprovalDialog, { PointsApprovalDialogInterface } from "../points/PointsApprovalDialog";
import CashoutList from "../points/CashoutList";
import PointRequestList from "../points/PointRequestList";
import PointsList from "../points/PointsList";
import SectionTitle from "../common/SectionTitle";

const ln = () => `[${moment().toISOString()}] ChildsPoints: `;

const ChildsPoints = () => {

    const approvalDialog = useRef<PointsApprovalDialogInterface>();
    const dispatch = useAppDispatch();

    const userId = useParams()["user_id"] || "none";
    const familyId = useParams()["family_id"] || "none";

    const store = useAppSelector((state) => ({
        points: state.points.userPoints,
        family: state.family.families[familyId],
        child: state.family.families[familyId].children[userId],
    }));

    const [loading, setLoading] = useState(true);

    const points = store.points;
    const child = store.child;

    const settledPoints = points.filter(x => x.status === PointStatus.SETTLED && x.request.type !== PointRequestType.CASHOUT);
    const requestedPoints = points.filter(x => x.status === PointStatus.WAITING);
    const cashouts = points.filter(x => x.status === PointStatus.SETTLED && x.request.type === PointRequestType.CASHOUT);

    const handleOnRequestClick = (p: PointSummary) => {
        const point = mapSummaryToLitePoint(p);
        point.user_id = userId;
        console.log("point", point);
        console.log("approvalDialog.current", approvalDialog.current);
        approvalDialog.current?.open(point, child);
    };

    useEffect(() => {
        setLoading(true);
        console.log(`${ln()}dispatching getUserPoints`);

        dispatch(getUserPoints(userId));
        setLoading(false);
    }, [dispatch]);

    return (
        <>
            {/* Left */}
            <div className="container mx-auto col-span-3">
                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500 mb-8">
                    <div className="card-body">
                        <SectionTitle>
                            {child.name}&apos;s points&nbsp;
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
                            {child.name}&apos;s requests&nbsp;
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
                            {child.name}&apos;s cashout history&nbsp;
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
