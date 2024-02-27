"use client"

import Image from "next/image";
import { ArrowDownRightIcon, ArrowUpRightIcon, CircleStackIcon, ClockIcon, CurrencyDollarIcon } from '@heroicons/react/24/solid'
import { AnimatedCounter } from  "react-animated-counter";
export default function Home() {
    return (
        <main className="overflow-hidden">
            {/* Top navbar */}
            {/* TK TK */}
        
            <section>
                <div className="w-screen xl:grid gap-8 xl:grid-cols-2 p-12">
                    {/* Left */}
                    <div className="container mx-auto p-4">
                        <div className="card soft-concave-shadow bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500 mb-8">
                            <figure className="px-10 pt-10">
                                <div className="hero-points-display rounded-xl p-16 bg-gradient-45 from-lime-500 to-lime-300 w-full text-center text-teal-600">
                                    <div className="text-9xl font-bold">
                                        213
                                    </div>
                                    <div className="text-7xl">
                                        points
                                    </div>
                                </div>
                            </figure>
                            <div className="card-body items-center text-center">
                                <p className="text-4xl">Wow, great job!</p>
                                <div className="card-actions mt-4">
                                    <button className="btn btn-primary sm:btn-lg">
                                        <CircleStackIcon className="w-8 h-8" />Request
                                    </button>
                                    <button className="btn btn-secondary sm:btn-lg">
                                        <CurrencyDollarIcon className="w-8 h-8" />Cashout
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Right */}
                    <div className="container mx-auto p-4">
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
        </main>
    );
}
