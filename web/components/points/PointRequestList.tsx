import React from "react";
import { ClockIcon } from "@heroicons/react/24/solid";

import { PointSummary } from "@/lib/models/Points";
import { formatDay_DDD_MMM_DD_hmm } from "@/lib/MomentUtils";

type Props = {
    points: PointSummary[];
    onClick?: (p: PointSummary) => void;
};

const PointRequestList = (props: Props) => {
    const { points } = props;

    const handleOnClick = (point: PointSummary) => {
        if (props.onClick)
            props.onClick(point);
    };

    return (
        <div className="list-container">
            {points.map((p, i) => {
                // const iColor = p.points > 0
                // ? "bg-green-400 text-green-700"
                // : "bg-red-400 text-red-700";

                // const pColor = p.points > 0
                //     ? "bg-green-500 text-green-100"
                //     : "bg-red-500 text-red-100";

                let wrapperClass = "list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200";
                if (props.onClick)
                    wrapperClass += " cursor-pointer";

                return <div key={i} className={wrapperClass} onClick={() => handleOnClick(p)}>
                    <div className="flex flex-row items-center space-x-4">
                        <ClockIcon className="bg-teal-400 text-teal-700 rounded-full p-2 h-12 w-12 " />
                        <div>
                            <p className="tracking-tight font-bold text-sm">{formatDay_DDD_MMM_DD_hmm(p.updated_on)}</p>
                            <p className="text-lg">{p.reason}</p>
                            {p.parent_notes
                                ? <p className="text-sm">{p.parent_notes}</p>
                                : <></>}
                        </div>
                    </div>
                    <div>
                        <div className={"bg-teal-500 text-teal-100 rounded-full px-4 py-1  text-lg font-bold"}>{p.points}</div>
                    </div>
                </div>;
            })}
        </div>
    );
};

export default PointRequestList;
