"use client";

import React, { useEffect, useState } from "react";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";

import { UserData } from "@/lib/auth/Auth";
import { MapType } from "@/lib/models/Common";
import { Family } from "@/lib/models/Family";
import { FamilySlice } from "@/slices/familySlice";
import { useAppDispatch } from "@/store/hooks";

import FamiliesList from "../family/FamiliesList";
import ChildsPoints from "../family/ChildsPoints";

type Props = {
    user: UserData;
    families: MapType<Family>;
};

const ParentsApp = (props: Props) => {
    // The following lines are needed as a workaround for document null error thrown by Next.
    const [render, setRender] = useState(false);
    useEffect(() => setRender(true), []);

    const dispatch = useAppDispatch();

    // Dispatch initial data
    dispatch(FamilySlice.actions.setFamilies(props.families));

    if (!render)
        return null;

    return (
        <section>
            <div className="w-screen gap-8 grid grid-cols-1 p-12">
                <div className="container mx-auto">
                    <Router>
                        <Routes>
                            <Route path="/app/parents" element={<FamiliesList />} />
                            <Route path="/app/parents/:family_id/points/:user_id" element={<ChildsPoints />} />
                        </Routes>
                    </Router>
                </div>                    
            </div>
        </section>
    );
};

export default ParentsApp;
