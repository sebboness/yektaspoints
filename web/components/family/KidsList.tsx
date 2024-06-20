"use client";

import Link from "next/link";
import { UserCircleIcon } from "@heroicons/react/24/solid";
import { useRouter } from "next/navigation";
import { useAppDispatch, useAppStore } from "@/store/hooks";
import { Family } from "@/lib/models/Family";
import { useState } from "react";

type Props = {
    initialFamily: Family;
}

const KidsList = ({ initialFamily }: Props) => {
    const router = useRouter();
    const dispatch = useAppDispatch();
    const store = useAppStore();

    const [family, setFamily] = useState<Family>(initialFamily);
    const [loading, setLoading] = useState(false);

    /**
     * Navigates to a family detail page for the given family ID
     * @param familyId The family ID
     */
    const goToFamily = (familyId: string) => {
        router.push(`/family/${familyId}`);
    };

    console.log("KidsList initialFamily", family);
    console.log("KidsList initialFamily children", family.children);

    Object.entries(family.children);

    return (
        <div className="list-container">

            {Object.entries(family.children).map((kvp, i) => {
                const child = kvp[1];

                return <div key={i} className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                    <div className="flex flex-row items-center space-x-4">
                        <UserCircleIcon className="bg-indigo-700 rounded-full p-2 h-12 w-12 text-indigo-200" />
                        <div>
                            <Link className="text-lg md:text-2xl" href={`/family/${family.family_id}/${child.user_id}`}>{child.name}</Link>
                        </div>
                    </div>
                </div>;
            })}
        </div>
    );
};

export default KidsList;
