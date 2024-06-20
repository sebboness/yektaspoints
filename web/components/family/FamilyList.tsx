"use client";

import React, { useEffect, useState } from "react";
import moment from "moment";
import { useRouter } from "next/navigation";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { UserCircleIcon } from "@heroicons/react/24/solid";

import { getUser } from "@/slices/authSlice";
import { useAppDispatch, useAppStore } from "@/store/hooks";

import CardSingleBody from "../common/CardSingleBody";
import SectionTitle from "../common/SectionTitle";

const ln = () => `[${moment().toISOString()}] FamilyList: `;

type Props = {
    initialFamilyIds: string[];
};

const FamilyList = ({ initialFamilyIds }: Props) => {
    const router = useRouter();
    const dispatch = useAppDispatch();
    const store = useAppStore();

    const [familyIds, setFamilyIds] = useState(initialFamilyIds);
    const [loading, setLoading] = useState(false);

    /**
     * Navigates to a family detail page for the given family ID
     * @param familyId The family ID
     */
    const goToFamily = (familyId: string) => {
        router.push(`/family/${familyId}`);
    };

    useEffect(() => {
        console.log(`${ln()}dispatching getUser`);

        setLoading(true);
        dispatch(getUser());
        setLoading(false);

        const user = store.getState().auth.user;
        if (user) {
            console.log(`${ln()}setting family IDs`);
            setFamilyIds(user.family_ids);
        };
    }, []);

    return (
        <CardSingleBody>
            <SectionTitle>
                My families&nbsp;
                {loading
                    ? <FontAwesomeIcon icon={faSpinner} spin />
                    : <></>}
            </SectionTitle>

            <div className="list-container">
                {familyIds.map((familyId, i) => (
                    <div
                        key={i}
                        className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200 cursor-pointer"
                        onClick={() => goToFamily(familyId)}>
                        <div className="flex flex-row items-center space-x-4">
                            <UserCircleIcon className="bg-indigo-700 rounded-full p-2 h-12 w-12 text-indigo-200" />
                            <div className="text-lg md:text-2xl">
                                Family {i+1}
                            </div>
                        </div>
                    </div>)
                )}
            </div>
        </CardSingleBody>
    );
};

export default FamilyList;
