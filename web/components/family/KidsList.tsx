"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import moment from "moment";
import { UserCircleIcon } from "@heroicons/react/24/solid";

import { getFamily } from "@/slices/familySlice";
import { useAppDispatch, useAppStore } from "@/store/hooks";
import { Family } from "@/lib/models/Family";
import CardSingleBody from "../common/CardSingleBody";
import SectionTitle from "../common/SectionTitle";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";

const ln = () => `[${moment().toISOString()}] KidsList: `;

type Props = {
    familyId: string
    initialFamily: Family;
}

const KidsList = ({ initialFamily, familyId }: Props) => {
    const router = useRouter();
    const dispatch = useAppDispatch();
    const store = useAppStore();

    const [family, setFamily] = useState<Family>(initialFamily);
    const [loading, setLoading] = useState(false);

    /**
     * Navigates to a family detail page for the given family ID
     * @param userId The family ID
     */
    const goToKid = (userId: string) => {
        router.push(`/family/${family.family_id}/${userId}`);
    };

    useEffect(() => {
        console.log(`${ln()}dispatching getFamily`);

        setLoading(true);
        dispatch(getFamily(familyId));
        setLoading(false);

        const fetchedFamily = store.getState().family.families[familyId];
        if (fetchedFamily) {
            console.log(`${ln()}setting family IDs`);
            setFamily(fetchedFamily);
        };
    }, []);

    const kidsKvp = Object.entries(family.children);

    return (
        <CardSingleBody>
            <SectionTitle>
                The little one{kidsKvp.length === 1 ? "" : "s"}&nbsp;
                {loading
                    ? <FontAwesomeIcon icon={faSpinner} spin />
                    : <></>}
            </SectionTitle>
            
            <div className="list-container">

                {kidsKvp.map((kvp, i) => {
                    const child = kvp[1];

                    return <div
                        key={i}
                        className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200 cursor-pointer"
                        title={`Manage ${child.name}'s points`}
                        onClick={() => goToKid(child.user_id)}>
                        <div className="flex flex-row items-center space-x-4">
                            <UserCircleIcon className="bg-indigo-700 rounded-full p-2 h-12 w-12 text-indigo-200" />
                            <div className="text-lg md:text-2xl">
                                {child.name}
                            </div>
                        </div>
                    </div>;
                })}
            </div>
        </CardSingleBody>
    );
};

export default KidsList;
