import { ArrowDownRightIcon, ArrowUpRightIcon } from "@heroicons/react/24/solid";

import React from "react";
import { UserPoints } from "@/lib/models/Points";
import { formatDay_DDD_MMM_DD_hmm } from "@/lib/MomentUtils";

type Props = {
    sum: UserPoints;
}

const RecentPoints = (props: Props) => {
    const { sum } = props;

    return (
        <div className="list-container">
            {sum.recent_points.map((p, i) => {
                const iColor = p.points > 0
                ? "bg-green-400 text-green-700"
                : "bg-red-400 text-red-700";

                const pColor = p.points > 0
                    ? "bg-green-500 text-green-100"
                    : "bg-red-500 text-red-100";

                return <div key={i} className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                    <div className="flex flex-row items-center space-x-4">
                        {p.points < 1
                            ? <ArrowDownRightIcon className={`${iColor} rounded-full p-2 h-12 w-12`} />
                            : <ArrowUpRightIcon className={`${iColor} rounded-full p-2 h-12 w-12`} />}
                        <div>
                            <p className="tracking-tight font-bold text-sm">{formatDay_DDD_MMM_DD_hmm(p.updated_on)}</p>
                            <p className="text-lg">{p.reason}</p>
                            {p.parent_notes
                                ? <p className=" text-sm">{p.parent_notes}</p>
                                : <></>}
                        </div>
                    </div>
                    <div>
                        <button className={`${pColor} rounded-full px-4 py-1 text-lg font-bold`}>{p.points}</button>
                    </div>
                </div>;
            })}
        </div>
    )
}

export default RecentPoints;