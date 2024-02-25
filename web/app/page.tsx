import Image from "next/image";

export default function Home() {
    return (
        <main className="overflow-hidden">
            {/* Top navbar */}
            <div className="navbar bg-base-100">
                <div className="flex-1">
                    <a className="btn btn-ghost text-xl">MyPoints</a>
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
                <div className="w-screen h-screen grid grid-rows-2 text-white text-4xl xl:grid-cols-2">
                    {/* Page 1 */}
                    <div className="w-full h-full bg-blue-800 centered md:h-screen">
                        <p>Page 1</p>
                    </div>

                    {/* Page 2 */}
                    <div className="w-full h-full bg-black centered md:h-screen">
                        <p>Page 2</p>
                    </div>
                </div>
            </section>
        </main>
    );
}
