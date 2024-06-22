/** @type {import('next').NextConfig} */
const nextConfig = {
    rewrites: async () => [
        {
          source: '/app/parents/:any*',
          destination: '/app/parents/',
        },
        {
          source: '/app/mypoints/:any*',
          destination: '/app/mypoints/',
        },
    ],
};

export default nextConfig;
