import { render, screen } from "@testing-library/react";
import Home from "@/app/page";
import { Providers } from "@/app/provider";

it('should have Wow! text', () => {
    render(<Providers><Home /></Providers>); // ARRANGE

    const myElem = screen.getByText('Wow!') // ACT

    expect(myElem).toBeInTheDocument(); // ASSERT
})