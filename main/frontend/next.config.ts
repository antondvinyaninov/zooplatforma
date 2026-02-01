import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ['@pet/shared'],
  turbopack: {},
  experimental: {
    allowedDevOrigins: [
      'localhost',
      'localhost:3000',
      'localhost:80',
      '127.0.0.1',
      '10.0.0.0/8', // Docker internal network
      'my-projects-zooplatforma.crv1ic.easypanel.host',
    ],
  },
};

export default nextConfig;
