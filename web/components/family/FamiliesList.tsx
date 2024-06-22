"use client";

import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import moment from "moment";
import { UserCircleIcon } from "@heroicons/react/24/solid";

import { useAppStore } from "@/store/hooks";

import CardSingleBody from "../common/CardSingleBody";

const ln = () => `[${moment().toISOString()}] FamiliesList: `;

type Props = {
    
};

const FamiliesList = (props: Props) => {
    const navigate = useNavigate();
    const store = useAppStore().getState();

    const [families, setFamilies] = useState(store.family.families);
    const familyIds = Object.entries(families);

    return familyIds.map((familyKvp, i) => {
        const family = familyKvp[1];
        const children = Object.entries(family.children);

        return <CardSingleBody key={i}>
            <div className="list-container">
                {children.map((childKvp, j) => {
                const child = childKvp[1];

                return <div
                    key={i}
                    className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200 cursor-pointer"
                    title={`Manage ${child.name}'s points`}
                    onClick={() => navigate(`/app/parents/${family.family_id}/points/${child.user_id}`)}>
                    <div className="flex flex-row items-center space-x-4">
                        <UserCircleIcon className="bg-indigo-700 rounded-full p-2 h-12 w-12 text-indigo-200" />
                        <div className="text-lg md:text-2xl">
                            {child.name}
                        </div>
                    </div>
                </div>;
            })}
            </div>
        </CardSingleBody>;
    });
};

export default FamiliesList;
