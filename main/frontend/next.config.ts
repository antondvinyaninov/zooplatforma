import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ['@pet/shared'],
  turbopack: {},
};

export default nextConfig;
