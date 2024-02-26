import Image from "next/image";

export default function Home() {
    return (
        <main className="overflow-hidden">
            {/* Top navbar */}
            <div className="navbar ">
                <div className="flex-1">
                    <a className="btn btn-ghost text-xl">My points</a>
                </div>
                <div className="flex-none">
                    <div className="dropdown dropdown-end">
                        <div tabIndex={0} role="button" className="btn btn-ghost btn-circle avatar">
                            <div className="w-10 rounded-full">
                                <img alt="Tailwind CSS Navbar component" src="https://daisyui.com/images/stock/photo-1534528741775-53994a69daeb.jpg" />
                            </div>
                        </div>
                        <ul tabIndex={0} className="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
                            <li>
                                <a className="justify-between">
                                    Profile
                                    <span className="badge">New</span>
                                </a>
                            </li>
                            <li><a>Logout</a></li>
                        </ul>
                    </div>
                </div>
            </div>
        
            <section>
                <div className="w-screen grid grid-rows-2 gap-8 xl:grid-cols-2 p-16">
                    {/* Left */}
                    <div className="container mx-auto p-4">
                        <div className="card soft-concave-shadow bg-gradient-135 from-pink-100 to-stone-50">
                            <figure className="px-10 pt-10">
                                <div className="hero-points-display rounded-xl p-16 bg-lime-400 w-full text-center text-teal-600">
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
                                <div className="card-actions">
                                    <button className="btn btn-primary btn-lg">Request points</button>
                                    <button className="btn btn-secondary btn-lg">Cashout points</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Right */}
                    <div className="container mx-auto p-4">
                        <div className="card soft-concave-shadow bg-gradient-135 from-pink-100 to-stone-50 mb-16">
                            <div className="card-body">
                                <div className="p-2">
                                    <p className="text-2xl font-bold">History</p>
                                </div>
                            </div>
                        </div>

                        <div className="card soft-concave-shadow bg-gradient-135 from-pink-100 to-stone-50">
                            <div className="card-body">
                                <p className="text-2xl font-bold">Open requests</p>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </main>
    );
}
