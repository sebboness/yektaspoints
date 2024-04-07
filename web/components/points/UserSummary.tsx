"use client";

import {
    CircleStackIcon,
    CurrencyDollarIcon,
    HandThumbUpIcon,
    LightBulbIcon,
    SparklesIcon,
} from "@heroicons/react/24/solid";
import React, { useEffect, useState } from "react";
import RequestPointsDialog, { requestPointsDialogID } from "./RequestPointsDialog";
import { useAppDispatch, useAppSelector } from "@/store/hooks";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Image from "next/image";
import PointRequestList from "./PointRequestList";
import PointsList from "./PointsList";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { getUserPointSummary } from "@/slices/pointsSlice";
import moment from "moment";

const ln = () => `[${moment().toISOString()}] UserSummary: `;

const UserSummary = () => {
    
    const [loading, setLoading] = useState(true);
    const dispatch = useAppDispatch();

    const user = useAppSelector((state) => state.auth.user);
    const sum = useAppSelector((state) => state.points.userSummary);
    const userID = user ? user.user_id : "";
    
    const openRequestPointsDialog = (e: React.MouseEvent<HTMLButtonElement>) => {
        if (typeof document !== "undefined") {
            const dialog = document.getElementById(requestPointsDialogID) as HTMLDialogElement;
            if (dialog) {
                dialog.showModal();
            }
        }

        e.preventDefault();
        return false;
    }

    useEffect(() => {
        if (userID) {
            console.log(`${ln()}dispatching getUserPointSummary`);
            dispatch(getUserPointSummary(userID));
            setLoading(false);
        }
    }, [userID]);

    console.log(`${ln()}info`, loading, userID, sum);

    if (!user)
        return <></>;

    const pointsGainedDisplay = `${sum.points_last_7_days} point${sum.points_last_7_days === 1 ? "" : "s"}`;
    const pointsLostDisplay = `${-sum.points_lost_last_7_days} point${sum.points_lost_last_7_days === -1 ? "" : "s"}`;

    return (
        <div className="w-screen xl:grid gap-8 xl:grid-cols-2 p-12">
            {/* Left */}
            <div className="container mx-auto">
                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500 mb-8">
                    <div className="card-body items-center">
                        <Image
                            className="md:mt-8"
                            src="/img/logo-512.svg"
                            width={768}
                            height={327}
                            priority={true}
                            alt="Picture of the author"
                        />

                        <div className="hero-points-display rounded-xl w-full text-center text-teal-600">    
                            <div className="text-7xl sm:text-9xl font-bold py-16">                        
                                {loading
                                    ? <><FontAwesomeIcon icon={faSpinner} spin /></>
                                    : <>{sum.balance}</>}
                            </div>
                        </div>
                    </div>
                    <div className="card-body">
                        <p className="text-xl sm:text-4xl"><b>Wow!</b> Great job, {user.name}!</p>
                        
                        {/* Info items */}
                        <div className="list-container">
                            {sum.points_last_7_days > 0
                                ? <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                    <div className="flex flex-row items-center space-x-4">
                                        <LightBulbIcon className="bg-orange-400 text-yellow-200 rounded-full p-2 h-12 w-12 " />
                                        <div>
                                            <p className="text-lg md:text-2xl">You earned {pointsGainedDisplay} this week. Keept it up!</p>
                                        </div>
                                    </div>
                                </div>
                                : <></>}

                            {sum.points_lost_last_7_days == 0
                                ? <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                    <div className="flex flex-row items-center space-x-4">
                                        <HandThumbUpIcon className="bg-blue-700 text-blue-300 rounded-full p-2 h-12 w-12 " />
                                        <div>
                                            <p className="text-lg md:text-2xl">You didn&apos;t lose any points this week!</p>
                                        </div>
                                    </div>
                                </div>
                                : <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                <div className="flex flex-row items-center space-x-4">
                                    <SparklesIcon className=" bg-gray-900 text-yellow-200 rounded-full p-2 h-12 w-12 " />
                                    <div>
                                        <p className="text-lg md:text-2xl">You lost {pointsLostDisplay} this week, but don't worry!</p>
                                        <p className="">You can earn more with good behavior and good grades.</p>
                                    </div>
                                </div>
                            </div>}
                        </div>
                        {/* End Info items */}
                    </div>

                    {/* Point actions */}
                    <div className="card-body items-center text-center">                                
                        <div className="card-actions mt-4">
                            <button className="btn btn-primary btn-lg" onClick={openRequestPointsDialog}>
                                <CircleStackIcon className="w-8 h-8" />Earn
                            </button>
                            <button className="btn btn-secondary btn-lg">
                                <CurrencyDollarIcon className="w-8 h-8" />Cashout
                            </button>
                        </div>
                        
                    </div>
                </div>
            </div>

            {/* Right */}
            <div className="container mx-auto">
                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 mb-16 border border-zinc-500">
                    <div className="card-body">
                        <p className="text-2xl font-bold">What I got so far</p>

                        {/* Recent points */}
                        <PointsList points={sum.recent_points} />
                        {/* End recent points */}

                        <button className="btn btn-primary btn-lg mt-4">See more...</button>
                    </div>
                </div>

                <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                    <div className="card-body">
                        <p className="text-2xl font-bold">Waiting on points</p>

                        {/* Points waiting */}
                        <PointRequestList points={sum.recent_requests} />
                        {/* End points waiting */}
                    </div>
                </div>
            </div>

            <RequestPointsDialog />
        </div>
    );
}

export default UserSummary;
