import React from "react";

type Props = {
    children?: React.ReactNode;
};

const SectionTitle = ({ children }: Props) => {
    return (
        <p className="text-2xl font-bold">{children}</p>
    );
};

export default SectionTitle;
