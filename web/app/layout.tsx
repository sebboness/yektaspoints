import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
    title: "My Points",
    description: "An app to help children learn the value of money",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode; }>) {
    return (
        <html lang="en" data-theme="cupcake">
        <body className={inter.className + " notebook"}>{children}</body>
        </html>
    );
}
