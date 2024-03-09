import React from "react";
import Image from "next/image";
import {
    ArrowDownRightIcon,
    ArrowUpRightIcon,
    CircleStackIcon,
    ClockIcon,
    CurrencyDollarIcon,
    HandThumbUpIcon,
    LightBulbIcon,
} from "@heroicons/react/24/solid";
import { AuthWrapper } from "@/components/AuthWrapper";

export default async function Home() {

    return (
        <AuthWrapper>
            {/* Top navbar */}
            {/* TK TK */}
        
            <section>
                <div className="w-screen xl:grid gap-8 xl:grid-cols-2 p-12">
                    {/* Left */}
                    <div className="container mx-auto">
                        <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500 mb-8">
                            <div className="card-body items-center">
                                <Image
                                    className="md:mt-8"
                                    src="/img/logo-512.svg"
                                    width={768}
                                    height={327}
                                    priority={true}
                                    alt="Picture of the author"
                                />

                                <div className="hero-points-display rounded-xl w-full text-center text-teal-600">
                                    <div className="text-7xl sm:text-9xl font-bold py-16">
                                        213
                                    </div>
                                </div>
                            </div>
                            <div className="card-body">
                                <p className="text-xl sm:text-4xl"><b>Wow!</b> Great job, Yekta!</p>
                                
                                {/* Info items */}
                                <div className="list-container">
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <LightBulbIcon className="bg-orange-400 rounded-full p-2 h-12 w-12 text-yellow-200" />
                                            <div>
                                                <h1 className="text-lg md:text-2xl">You earned 47 points this week. Keept it up!</h1>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <HandThumbUpIcon className="bg-blue-700 rounded-full p-2 h-12 w-12 text-blue-300" />
                                            <div>
                                                <h1 className="text-lg md:text-2xl">You didn&apos;t lose any points this week!</h1>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                {/* End Info items */}
                            </div>

                            {/* Point actions */}
                            <div className="card-body items-center text-center">                                
                                <div className="card-actions mt-4">
                                    <button className="btn btn-primary btn-lg">
                                        <CircleStackIcon className="w-8 h-8" />Earn
                                    </button>
                                    <button className="btn btn-secondary btn-lg">
                                        <CurrencyDollarIcon className="w-8 h-8" />Cashout
                                    </button>
                                </div>
                                
                            </div>
                        </div>
                    </div>

                    {/* Right */}
                    <div className="container mx-auto">
                        <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 mb-16 border border-zinc-500">
                            <div className="card-body">
                                <p className="text-2xl font-bold">What I got so far</p>

                                {/* Items */}
                                <div className="list-container">
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <ArrowUpRightIcon className="bg-green-400 rounded-full p-2 h-12 w-12 text-green-700" />
                                            <div>
                                                <h1 className="tracking-tight">Mon, Feb 26th</h1>
                                                <p className="font-light">Cleaning my room</p>
                                            </div>
                                        </div>
                                        <div>
                                            <button className="bg-green-500 rounded-full px-4 py-1 text-green-100 text-lg font-bold">10</button>
                                        </div>
                                    </div>
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <ArrowDownRightIcon className="bg-red-400 rounded-full p-2 h-12 w-12 text-red-700" />
                                            <div>
                                                <h1 className="tracking-tight">Sun, Feb 25th</h1>
                                                <p className="font-light">Nagging and not listening</p>
                                            </div>
                                        </div>
                                        <div>
                                            <button className="bg-red-500 rounded-full px-4 py-1 text-red-100 text-lg font-bold">5</button>
                                        </div>
                                    </div>
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <ArrowUpRightIcon className="bg-green-400 rounded-full p-2 h-12 w-12 text-green-700" />
                                            <div>
                                                <h1 className="tracking-tight">Fri, Feb 23rd</h1>
                                                <p className="font-light">Taking out Gina for a walk</p>
                                            </div>
                                        </div>
                                        <div>
                                            <button className="bg-green-500 rounded-full px-4 py-1 text-green-100 text-lg font-bold">2</button>
                                        </div>
                                    </div>
                                </div>
                                {/* End Items */}

                                <button className="btn btn-primary btn-lg mt-4">See more...</button>
                            </div>
                        </div>

                        <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                            <div className="card-body">
                                <p className="text-2xl font-bold">Waiting on points</p>
                                
                                {/* Items */}
                                <div className="list-container">
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <ClockIcon className="bg-teal-400 rounded-full p-2 h-12 w-12 text-teal-700" />
                                            <div>
                                                <h1 className="tracking-tight">Mon, Feb 26th</h1>
                                                <p className="font-light">Cleaning my room</p>
                                            </div>
                                        </div>
                                        <div>
                                            <button className="bg-teal-500 rounded-full px-4 py-1 text-teal-100 text-lg font-bold">10</button>
                                        </div>
                                    </div>
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <ClockIcon className="bg-teal-400 rounded-full p-2 h-12 w-12 text-teal-700" />
                                            <div>
                                                <h1 className="tracking-tight">Sun, Feb 25th</h1>
                                                <p className="font-light">Nagging and not listening</p>
                                            </div>
                                        </div>
                                        <div>
                                            <button className="bg-teal-500 rounded-full px-4 py-1 text-teal-100 text-lg font-bold">5</button>
                                        </div>
                                    </div>
                                    <div className="list-items flex flex-row items-center justify-between mx-auto border-4 border-zinc-500 py-4 rounded-full my-4 px-4 bg-gradient-135 from-base-100 to-base-200">
                                        <div className="flex flex-row items-center space-x-4">
                                            <ClockIcon className="bg-teal-400 rounded-full p-2 h-12 w-12 text-teal-700" />
                                            <div>
                                                <h1 className="tracking-tight">Fri, Feb 23rd</h1>
                                                <p className="font-light">Taking out Gina for a walk</p>
                                            </div>
                                        </div>
                                        <div>
                                            <button className="bg-teal-500 rounded-full px-4 py-1 text-teal-100 text-lg font-bold">2</button>
                                        </div>
                                    </div>
                                </div>
                                {/* End Items */}
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </AuthWrapper>
        );
}
