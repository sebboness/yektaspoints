import React from "react";

const CardShell = ({ children }: { children: React.ReactNode; }) => {
    return (
        <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500 mb-8">
            {children}
        </div>
    );
}

export default CardShell;
