import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@pet/shared'],
  
  // Пустой turbopack config для совместимости с Next.js 16
  turbopack: {},
  
  typescript: {
    ignoreBuildErrors: process.env.NODE_ENV === 'development',
  },
};

export default nextConfig;
