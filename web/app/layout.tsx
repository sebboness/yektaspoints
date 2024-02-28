import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "../styles/globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
    title: "My Points",
    description: "An app to help children learn the value of money",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode; }>) {
    return (
        <html lang="en" data-theme="cupcake">
        <head>
            <link rel="apple-touch-icon" sizes="180x180" href="/img/favicon-180.png" />
            <link rel="icon" href="/img/favicon.ico" />
            <link rel="icon" type="image/png" sizes="32x32" href="/img/favicon-32.png" />
            <link rel="icon" type="image/png" sizes="96x96" href="/img/favicon-96.png" />
            <link rel="icon" type="image/png" sizes="192x192" href="/img/favicon-192.png" />

            <meta name="msapplication-TileColor" content="#00a7ff" /> 
            <meta name="msapplication-TileImage" content="" />
            <meta name="theme-color" content="#00a7ff" />

            <meta name="description" content="TK TK TK" />
        
            <meta property="og:type" content="website" />
            <meta property="og:url" content="https://mypoints.hexonite.net/" />
            <meta property="og:title" content="myPoints - by Hexonite.net" />
            <meta property="og:description" content="TK TK TK" />
            <meta property="og:image" content="https://mypoints.hexonite.net/img/landing-image.jpg" />

            <meta property="twitter:card" content="summary_large_image" />
            <meta property="twitter:url" content="https://mypoints.hexonite.net/" />
            <meta property="twitter:title" content="myPoints - by Hexonite.net" />
            <meta property="twitter:description" content="TK TK TK" />
            <meta property="twitter:image" content="https://mypoints.hexonite.net/img/landing-image.jpg" />
            
            <meta name="viewport" content="width=device-width, initial-scale=1" />
            <meta httpEquiv="x-ua-compatible" content="ie=edge" />
            <meta name="robots" content="noindex, follow" />
        </head>
        <body className={inter.className + " notebook"}>
            <main className="overflow-hidden font-primary">
                {children}
            </main>
        </body>
        </html>
    );
}
