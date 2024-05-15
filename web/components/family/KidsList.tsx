"use client";

import React, { useState } from "react";
import Link from "next/link";
import { UserCircleIcon } from "@heroicons/react/24/solid";

type Kid = {
    id: string;
    name: string;
};

const initialKids: Kid[] = [
    {
        id: "123",
        name: "Yekta",
    },
];

const KidsList = () => {
    // const user = useAppSelector((state) => state.auth.user);
    const [kids, setKids] = useState(initialKids);

    console.log("setKids", setKids);

    return (
        <div className="list-container">
            {kids.map((x, i) => (
                <div key={i} className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                    <div className="flex flex-row items-center space-x-4">
                        <UserCircleIcon className="bg-indigo-700 rounded-full p-2 h-12 w-12 text-indigo-200" />
                        <div>
                            <Link className="text-lg md:text-2xl" href={`/family/${x.id}`}>{x.name}</Link>
                        </div>
                    </div>
                </div>)
            )}
        </div>
    );
};

export default KidsList;
