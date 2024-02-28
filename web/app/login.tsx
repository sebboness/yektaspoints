import Image from "next/image";

export default function Home() {
    return (
        <>
            {/* Top navbar */}
            {/* TK TK */}
        
            <div className="hero min-h-screen">
                <div className="card card-compact bg-base-100 shadow-xl p-4 bg-gradient-135 from-pink-200 to-lime-100 border border-zinc-500">
                    <form className="card-body">

                        <figure className="mb-8">
                            <Image
                                src="/img/logo-512.svg"
                                width={384}
                                height={164}
                                priority={true}
                                alt="Picture of the author"
                                />
                        </figure>

                        <h1 className="text-5xl">Hello there! :)</h1>

                        <div className="form-control">
                            <input type="email" placeholder="email" className="input input-bordered" required />
                        </div>
                        <div className="form-control">
                            <input type="password" placeholder="password" className="input input-bordered" required />
                            <label className="label">
                                <a href="#" className="label-text-alt link link-hover">Forgot password?</a>
                            </label>
                        </div>
                        <div className="form-control mt-6">
                            <button className="btn btn-primary">Login</button>
                        </div>
                    </form>
                </div>
            </div>
        </>
    );
}
