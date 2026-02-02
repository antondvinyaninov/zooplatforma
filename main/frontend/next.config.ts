import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ['@pet/shared'],
  turbopack: {},
  async rewrites() {
    // В production используем относительные пути (backend на том же хосте)
    // В development используем localhost
    const apiUrl = process.env.NODE_ENV === 'production' 
      ? 'http://localhost:8000'  // В Docker контейнере backend на localhost:8000
      : 'http://localhost:8000';
    
    return [
      {
        source: '/api/:path*',
        destination: `${apiUrl}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
