/** @type {import('next').NextConfig} */
const nextConfig = {
    output: "standalone",
    transpilePackages: ["@nextui-org/react"],
    productionBrowserSourceMaps: false,
    experimental: {
        serverSourceMaps: false,
    },
    // Add any other existing configurations here
};

export default nextConfig;
