import type { NextConfig } from "next";

const apiTarget = process.env.API_INTERNAL_URL ?? "http://127.0.0.1:8080";

const nextConfig: NextConfig = {
  cacheComponents: true,
  partialPrefetching: true,
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: `${apiTarget}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;