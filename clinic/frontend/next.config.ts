import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  turbopack: {},
  async rewrites() {
    return [
      {
        source: '/api/petbase/:path*',
        destination: 'http://localhost:8100/api/:path*',
      },
      {
        source: '/api/main/:path*',
        destination: 'http://localhost:8000/api/:path*',
      },
    ];
  },
};

export default nextConfig;
